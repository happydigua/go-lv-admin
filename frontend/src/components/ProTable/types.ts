import type { DataTableColumns } from 'naive-ui';

export interface ProTableProps {
    request: (params: any) => Promise<{ list: any[]; total: number }>;
    columns: DataTableColumns<any>;
    searchSchema?: Record<string, SearchColumn>; // key: fieldName
    title?: string;
    rowKey?: string | ((row: any) => string | number);
}

export interface SearchColumn {
    type: 'input' | 'select' | 'date-picker';
    label: string;
    placeholder?: string;
    options?: { label: string; value: string | number }[]; // For Select
    defaultValue?: any;
}
