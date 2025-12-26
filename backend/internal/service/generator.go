package service

import (
	"bytes"
	"fmt"
	"go-lv-vue-admin/internal/global"
	"strings"
	"text/template"
	"unicode"
)

type GeneratorService struct{}

// TableInfo 表信息
type TableInfo struct {
	TableName    string `json:"tableName"`
	TableComment string `json:"tableComment"`
}

// ColumnInfo 列信息
type ColumnInfo struct {
	ColumnName    string `json:"columnName"`
	DataType      string `json:"dataType"`
	ColumnComment string `json:"columnComment"`
	IsNullable    string `json:"isNullable"`
	ColumnKey     string `json:"columnKey"`
	Extra         string `json:"extra"`
	// 生成配置
	GoField   string `json:"goField"`
	GoType    string `json:"goType"`
	JsonField string `json:"jsonField"`
	FormType  string `json:"formType"`  // input, select, textarea, date, etc.
	QueryType string `json:"queryType"` // eq, like, between, etc.
	IsQuery   bool   `json:"isQuery"`
	IsList    bool   `json:"isList"`
	IsForm    bool   `json:"isForm"`
}

// GenerateConfig 生成配置
type GenerateConfig struct {
	TableName    string       `json:"tableName"`
	TableComment string       `json:"tableComment"`
	ModuleName   string       `json:"moduleName"`  // 模块名，如 article
	PackageName  string       `json:"packageName"` // 包名，如 blog
	StructName   string       `json:"structName"`  // 结构体名，如 Article
	Columns      []ColumnInfo `json:"columns"`
}

// GetTables 获取数据库所有表
func (s *GeneratorService) GetTables() ([]TableInfo, error) {
	var tables []TableInfo

	dbName := extractDbName(global.LV_CONFIG.Database.Source)
	sql := `SELECT table_name, table_comment FROM information_schema.tables 
			WHERE table_schema = ? AND table_type = 'BASE TABLE'
			ORDER BY table_name`

	rows, err := global.LV_DB.Raw(sql, dbName).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t TableInfo
		if err := rows.Scan(&t.TableName, &t.TableComment); err != nil {
			continue
		}
		tables = append(tables, t)
	}
	return tables, nil
}

// GetTableColumns 获取表的列信息
func (s *GeneratorService) GetTableColumns(tableName string) ([]ColumnInfo, error) {
	var columns []ColumnInfo

	dbName := extractDbName(global.LV_CONFIG.Database.Source)
	sql := `SELECT column_name, data_type, column_comment, is_nullable, column_key, extra
			FROM information_schema.columns 
			WHERE table_schema = ? AND table_name = ?
			ORDER BY ordinal_position`

	rows, err := global.LV_DB.Raw(sql, dbName, tableName).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c ColumnInfo
		if err := rows.Scan(&c.ColumnName, &c.DataType, &c.ColumnComment, &c.IsNullable, &c.ColumnKey, &c.Extra); err != nil {
			continue
		}
		// 自动推断字段配置
		c.GoField = toCamelCase(c.ColumnName, true)
		c.GoType = dbTypeToGoType(c.DataType)
		c.JsonField = toCamelCase(c.ColumnName, false)
		c.FormType = inferFormType(c.DataType, c.ColumnName)
		c.QueryType = "eq"
		c.IsQuery = isQueryField(c.ColumnName)
		c.IsList = true
		c.IsForm = !isAutoField(c.ColumnName)
		columns = append(columns, c)
	}
	return columns, nil
}

// GenerateCode 生成代码
func (s *GeneratorService) GenerateCode(config GenerateConfig) (map[string]string, error) {
	result := make(map[string]string)

	// 生成 Model
	modelCode, err := s.generateModel(config)
	if err != nil {
		return nil, err
	}
	result["model"] = modelCode

	// 生成 Service
	serviceCode, err := s.generateService(config)
	if err != nil {
		return nil, err
	}
	result["service"] = serviceCode

	// 生成 API
	apiCode, err := s.generateAPI(config)
	if err != nil {
		return nil, err
	}
	result["api"] = apiCode

	// 生成 Router
	routerCode, err := s.generateRouter(config)
	if err != nil {
		return nil, err
	}
	result["router"] = routerCode

	// 生成前端 Vue 页面
	vueCode, err := s.generateVue(config)
	if err != nil {
		return nil, err
	}
	result["vue"] = vueCode

	// 生成前端 API
	frontendAPICode, err := s.generateFrontendAPI(config)
	if err != nil {
		return nil, err
	}
	result["frontendApi"] = frontendAPICode

	return result, nil
}

// 生成 Model
func (s *GeneratorService) generateModel(config GenerateConfig) (string, error) {
	tmpl := `package model

import "gorm.io/gorm"

// {{.StructName}} {{.TableComment}}
type {{.StructName}} struct {
	gorm.Model
{{- range .Columns}}
{{- if not (isAutoField .ColumnName)}}
	{{.GoField}} {{.GoType}} ` + "`" + `json:"{{.JsonField}}" gorm:"{{gormTag .}}"` + "`" + `{{if .ColumnComment}} // {{.ColumnComment}}{{end}}
{{- end}}
{{- end}}
}

func ({{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
`
	return s.executeTemplate(tmpl, config)
}

// 生成 Service
func (s *GeneratorService) generateService(config GenerateConfig) (string, error) {
	tmpl := `package service

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
)

type {{.StructName}}Service struct{}

// GetList 获取{{.TableComment}}列表
func (s *{{.StructName}}Service) GetList(page, pageSize int{{range .Columns}}{{if .IsQuery}}, {{.JsonField}} {{.GoType}}{{end}}{{end}}) ([]model.{{.StructName}}, int64, error) {
	var list []model.{{.StructName}}
	var total int64

	db := global.LV_DB.Model(&model.{{.StructName}}{})
{{range .Columns}}{{if .IsQuery}}
	if {{.JsonField}} != {{zeroValue .GoType}} {
		db = db.Where("{{.ColumnName}} {{queryOp .QueryType}} ?", {{queryValue .}})
	}
{{end}}{{end}}
	db.Count(&total)

	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Find(&list).Error

	return list, total, err
}

// GetById 根据ID获取{{.TableComment}}
func (s *{{.StructName}}Service) GetById(id uint) (*model.{{.StructName}}, error) {
	var item model.{{.StructName}}
	err := global.LV_DB.First(&item, id).Error
	return &item, err
}

// Create 创建{{.TableComment}}
func (s *{{.StructName}}Service) Create(item *model.{{.StructName}}) error {
	return global.LV_DB.Create(item).Error
}

// Update 更新{{.TableComment}}
func (s *{{.StructName}}Service) Update(item *model.{{.StructName}}) error {
	return global.LV_DB.Model(item).Updates(item).Error
}

// Delete 删除{{.TableComment}}
func (s *{{.StructName}}Service) Delete(id uint) error {
	return global.LV_DB.Delete(&model.{{.StructName}}{}, id).Error
}
`
	return s.executeTemplate(tmpl, config)
}

// 生成 API
func (s *GeneratorService) generateAPI(config GenerateConfig) (string, error) {
	tmpl := `package v1

import (
	"go-lv-vue-admin/internal/global"
	"go-lv-vue-admin/internal/model"
	"go-lv-vue-admin/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type {{.StructName}}Api struct{}

var {{.ModuleName}}Service = service.{{.StructName}}Service{}

// GetList 获取{{.TableComment}}列表
func (a *{{.StructName}}Api) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
{{range .Columns}}{{if .IsQuery}}
	{{.JsonField}} := c.Query("{{.JsonField}}")
{{end}}{{end}}

	list, total, err := {{.ModuleName}}Service.GetList(page, pageSize{{range .Columns}}{{if .IsQuery}}, {{convertQueryParam .}}{{end}}{{end}})
	if err != nil {
		global.LV_LOG.Error("获取{{.TableComment}}列表失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "获取失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{"list": list, "total": total, "page": page, "pageSize": pageSize},
		"msg":  "success",
	})
}

// GetById 获取{{.TableComment}}详情
func (a *{{.StructName}}Api) GetById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	item, err := {{.ModuleName}}Service.GetById(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"code": 7, "msg": "获取失败"})
		return
	}
	c.JSON(200, gin.H{"code": 0, "data": item, "msg": "success"})
}

// Create 创建{{.TableComment}}
func (a *{{.StructName}}Api) Create(c *gin.Context) {
	var item model.{{.StructName}}
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": err.Error()})
		return
	}

	if err := {{.ModuleName}}Service.Create(&item); err != nil {
		global.LV_LOG.Error("创建{{.TableComment}}失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "创建失败"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "创建成功"})
}

// Update 更新{{.TableComment}}
func (a *{{.StructName}}Api) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var item model.{{.StructName}}
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(400, gin.H{"code": 7, "msg": err.Error()})
		return
	}
	item.ID = uint(id)

	if err := {{.ModuleName}}Service.Update(&item); err != nil {
		global.LV_LOG.Error("更新{{.TableComment}}失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "更新失败"})
		return
	}

	c.JSON(200, gin.H{"code": 0, "msg": "更新成功"})
}

// Delete 删除{{.TableComment}}
func (a *{{.StructName}}Api) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := {{.ModuleName}}Service.Delete(uint(id)); err != nil {
		global.LV_LOG.Error("删除{{.TableComment}}失败", zap.Error(err))
		c.JSON(500, gin.H{"code": 7, "msg": "删除失败"})
		return
	}
	c.JSON(200, gin.H{"code": 0, "msg": "删除成功"})
}
`
	return s.executeTemplate(tmpl, config)
}

// 生成 Router
func (s *GeneratorService) generateRouter(config GenerateConfig) (string, error) {
	tmpl := `// {{.TableComment}}路由 - 添加到 router.go 的 privateGroup 中
{{.ModuleName}}Api := v1.{{.StructName}}Api{}
{{.ModuleName}}Group := privateGroup.Group("{{.PackageName}}/{{.ModuleName}}")
{
	{{.ModuleName}}Group.GET("list", {{.ModuleName}}Api.GetList)
	{{.ModuleName}}Group.GET(":id", {{.ModuleName}}Api.GetById)
	{{.ModuleName}}Group.POST("", {{.ModuleName}}Api.Create)
	{{.ModuleName}}Group.PUT(":id", {{.ModuleName}}Api.Update)
	{{.ModuleName}}Group.DELETE(":id", {{.ModuleName}}Api.Delete)
}
`
	return s.executeTemplate(tmpl, config)
}

// 生成前端 Vue 页面
func (s *GeneratorService) generateVue(config GenerateConfig) (string, error) {
	tmpl := `<template>
  <n-card title="{{.TableComment}}">
    <template #header-extra>
      <n-button type="primary" @click="handleAdd">新增</n-button>
    </template>

    <!-- 搜索区域 -->
    <n-space style="margin-bottom: 16px;">
{{- range .Columns}}{{if .IsQuery}}
      <n-input v-model:value="searchForm.{{.JsonField}}" placeholder="{{.ColumnComment}}" clearable style="width: 150px;" />
{{- end}}{{end}}
      <n-button type="primary" @click="fetchData">搜索</n-button>
      <n-button @click="handleReset">重置</n-button>
    </n-space>

    <n-data-table
      :columns="columns"
      :data="tableData"
      :pagination="pagination"
      :loading="loading"
      :bordered="false"
      @update:page="handlePageChange"
      @update:page-size="handlePageSizeChange"
    />
  </n-card>

  <!-- 编辑弹窗 -->
  <n-modal v-model:show="showModal" preset="dialog" :title="modalTitle" style="width: 600px;">
    <n-form ref="formRef" :model="formData" :rules="formRules" label-placement="left" label-width="80">
{{- range .Columns}}{{if .IsForm}}
      <n-form-item label="{{.ColumnComment}}" path="{{.JsonField}}">
{{- if eq .FormType "textarea"}}
        <n-input v-model:value="formData.{{.JsonField}}" type="textarea" placeholder="请输入{{.ColumnComment}}" />
{{- else if eq .FormType "number"}}
        <n-input-number v-model:value="formData.{{.JsonField}}" style="width: 100%;" />
{{- else if eq .FormType "date"}}
        <n-date-picker v-model:value="formData.{{.JsonField}}" type="date" style="width: 100%;" />
{{- else}}
        <n-input v-model:value="formData.{{.JsonField}}" placeholder="请输入{{.ColumnComment}}" />
{{- end}}
      </n-form-item>
{{- end}}{{end}}
    </n-form>
    <template #action>
      <n-button @click="showModal = false">取消</n-button>
      <n-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</n-button>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { h, ref, reactive, onMounted } from 'vue';
import { NButton, NSpace, useMessage, useDialog } from 'naive-ui';
import { get{{.StructName}}List, create{{.StructName}}, update{{.StructName}}, delete{{.StructName}} } from '@/api/{{.PackageName}}/{{.ModuleName}}';

const message = useMessage();
const dialog = useDialog();

const loading = ref(false);
const submitLoading = ref(false);
const showModal = ref(false);
const isEdit = ref(false);
const modalTitle = ref('新增{{.TableComment}}');
const formRef = ref();
const tableData = ref<any[]>([]);

const searchForm = reactive({ {{range .Columns}}{{if .IsQuery}}{{.JsonField}}: '',{{end}}{{end}} });

const pagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50]
});

const formData = ref<any>({ {{range .Columns}}{{if .IsForm}}{{.JsonField}}: {{defaultValue .GoType}},{{end}}{{end}} });

const formRules = { {{range .Columns}}{{if and .IsForm (eq .IsNullable "NO")}}
  {{.JsonField}}: { required: true, message: '请输入{{.ColumnComment}}', trigger: 'blur' },{{end}}{{end}}
};

const columns = [
  { title: 'ID', key: 'ID', width: 80 },
{{- range .Columns}}{{if .IsList}}
  { title: '{{.ColumnComment}}', key: '{{.JsonField}}' },
{{- end}}{{end}}
  {
    title: '操作',
    key: 'actions',
    width: 180,
    render: (row: any) => h(NSpace, null, {
      default: () => [
        h(NButton, { size: 'small', tertiary: true, type: 'info', onClick: () => handleEdit(row) }, { default: () => '编辑' }),
        h(NButton, { size: 'small', tertiary: true, type: 'error', onClick: () => handleDelete(row) }, { default: () => '删除' })
      ]
    })
  }
];

const fetchData = async () => {
  loading.value = true;
  try {
    const res: any = await get{{.StructName}}List({
      page: pagination.page,
      pageSize: pagination.pageSize,
      ...searchForm
    });
    tableData.value = res.list || [];
    pagination.itemCount = res.total || 0;
  } catch (error) {
    console.error('Failed to fetch data:', error);
  } finally {
    loading.value = false;
  }
};

const handlePageChange = (page: number) => { pagination.page = page; fetchData(); };
const handlePageSizeChange = (pageSize: number) => { pagination.pageSize = pageSize; pagination.page = 1; fetchData(); };
const handleReset = () => { Object.keys(searchForm).forEach(k => (searchForm as any)[k] = ''); pagination.page = 1; fetchData(); };

const handleAdd = () => {
  isEdit.value = false;
  modalTitle.value = '新增{{.TableComment}}';
  formData.value = { {{range .Columns}}{{if .IsForm}}{{.JsonField}}: {{defaultValue .GoType}},{{end}}{{end}} };
  showModal.value = true;
};

const handleEdit = (row: any) => {
  isEdit.value = true;
  modalTitle.value = '编辑{{.TableComment}}';
  formData.value = { ...row };
  showModal.value = true;
};

const handleDelete = (row: any) => {
  dialog.error({
    title: '删除确认',
    content: '确定要删除该{{.TableComment}}吗？',
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await delete{{.StructName}}(row.ID);
        message.success('删除成功');
        fetchData();
      } catch (error) {
        message.error('删除失败');
      }
    }
  });
};

const handleSubmit = () => {
  formRef.value?.validate(async (errors: any) => {
    if (!errors) {
      submitLoading.value = true;
      try {
        if (isEdit.value) {
          await update{{.StructName}}(formData.value.ID, formData.value);
          message.success('更新成功');
        } else {
          await create{{.StructName}}(formData.value);
          message.success('创建成功');
        }
        showModal.value = false;
        fetchData();
      } catch (error) {
        message.error(isEdit.value ? '更新失败' : '创建失败');
      } finally {
        submitLoading.value = false;
      }
    }
  });
};

onMounted(() => { fetchData(); });
</script>
`
	return s.executeTemplate(tmpl, config)
}

// 生成前端 API
func (s *GeneratorService) generateFrontendAPI(config GenerateConfig) (string, error) {
	tmpl := `import request from '@/utils/request';

// 获取{{.TableComment}}列表
export const get{{.StructName}}List = (params: any) => {
  return request({ url: '/{{.PackageName}}/{{.ModuleName}}/list', method: 'get', params });
};

// 获取{{.TableComment}}详情
export const get{{.StructName}}ById = (id: number) => {
  return request({ url: ` + "`" + `/{{.PackageName}}/{{.ModuleName}}/${id}` + "`" + `, method: 'get' });
};

// 创建{{.TableComment}}
export const create{{.StructName}} = (data: any) => {
  return request({ url: '/{{.PackageName}}/{{.ModuleName}}', method: 'post', data });
};

// 更新{{.TableComment}}
export const update{{.StructName}} = (id: number, data: any) => {
  return request({ url: ` + "`" + `/{{.PackageName}}/{{.ModuleName}}/${id}` + "`" + `, method: 'put', data });
};

// 删除{{.TableComment}}
export const delete{{.StructName}} = (id: number) => {
  return request({ url: ` + "`" + `/{{.PackageName}}/{{.ModuleName}}/${id}` + "`" + `, method: 'delete' });
};
`
	return s.executeTemplate(tmpl, config)
}

// 执行模板
func (s *GeneratorService) executeTemplate(tmplStr string, config GenerateConfig) (string, error) {
	funcMap := template.FuncMap{
		"isAutoField":       isAutoField,
		"gormTag":           gormTag,
		"zeroValue":         zeroValue,
		"queryOp":           queryOp,
		"queryValue":        queryValue,
		"convertQueryParam": convertQueryParam,
		"defaultValue":      defaultValue,
	}

	tmpl, err := template.New("gen").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, config); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// 辅助函数
func toCamelCase(s string, upperFirst bool) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if i == 0 && !upperFirst {
			parts[i] = strings.ToLower(part)
		} else {
			parts[i] = strings.Title(part)
		}
	}
	return strings.Join(parts, "")
}

func dbTypeToGoType(dbType string) string {
	dbType = strings.ToLower(dbType)
	switch {
	case strings.Contains(dbType, "int"):
		if strings.Contains(dbType, "bigint") {
			return "int64"
		}
		return "int"
	case strings.Contains(dbType, "float"), strings.Contains(dbType, "double"), strings.Contains(dbType, "decimal"):
		return "float64"
	case strings.Contains(dbType, "bool"):
		return "bool"
	case strings.Contains(dbType, "datetime"), strings.Contains(dbType, "timestamp"):
		return "time.Time"
	case strings.Contains(dbType, "date"):
		return "time.Time"
	default:
		return "string"
	}
}

func inferFormType(dbType, columnName string) string {
	if strings.Contains(columnName, "content") || strings.Contains(columnName, "desc") || strings.Contains(columnName, "remark") {
		return "textarea"
	}
	if strings.Contains(dbType, "int") || strings.Contains(dbType, "float") || strings.Contains(dbType, "decimal") {
		return "number"
	}
	if strings.Contains(dbType, "date") || strings.Contains(dbType, "time") {
		return "date"
	}
	return "input"
}

func isQueryField(columnName string) bool {
	return columnName == "name" || columnName == "title" || columnName == "status" || strings.Contains(columnName, "name")
}

func isAutoField(columnName string) bool {
	auto := []string{"id", "created_at", "updated_at", "deleted_at"}
	for _, a := range auto {
		if columnName == a {
			return true
		}
	}
	return false
}

func gormTag(c ColumnInfo) string {
	var tags []string
	if c.ColumnComment != "" {
		tags = append(tags, fmt.Sprintf("comment:%s", c.ColumnComment))
	}
	if c.ColumnKey == "UNI" {
		tags = append(tags, "unique")
	}
	return strings.Join(tags, ";")
}

func zeroValue(goType string) string {
	switch goType {
	case "int", "int64", "float64":
		return "0"
	case "bool":
		return "false"
	default:
		return `""`
	}
}

func queryOp(queryType string) string {
	switch queryType {
	case "like":
		return "LIKE"
	case "gt":
		return ">"
	case "lt":
		return "<"
	case "gte":
		return ">="
	case "lte":
		return "<="
	default:
		return "="
	}
}

func queryValue(c ColumnInfo) string {
	if c.QueryType == "like" {
		return `"%"+` + c.JsonField + `+"%"`
	}
	return c.JsonField
}

func convertQueryParam(c ColumnInfo) string {
	if c.GoType == "int" || c.GoType == "int64" {
		return fmt.Sprintf("parseInt(%s)", c.JsonField)
	}
	return c.JsonField
}

func defaultValue(goType string) string {
	switch goType {
	case "int", "int64", "float64":
		return "0"
	case "bool":
		return "false"
	default:
		return `''`
	}
}

func parseInt(s string) string {
	return fmt.Sprintf("parseInt(%s, 10)", s)
}

func firstLower(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// extractDbName 从 DSN 中解析数据库名
// 假设 DSN 格式: user:password@tcp(host:port)/dbname?params
func extractDbName(dsn string) string {
	// 查找 "/" 后的数据库名
	slashIndex := strings.LastIndex(dsn, "/")
	if slashIndex == -1 {
		return ""
	}
	dbPart := dsn[slashIndex+1:]
	// 去除参数
	if questionIndex := strings.Index(dbPart, "?"); questionIndex != -1 {
		dbPart = dbPart[:questionIndex]
	}
	return dbPart
}
