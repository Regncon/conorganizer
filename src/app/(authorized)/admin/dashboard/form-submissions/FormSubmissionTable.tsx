'use client';
import {
    DataGrid,
    useGridApiRef,
    type GridEventListener,
    type GridRowParams,
    GridToolbar,
    type GridRowId,
    type GridCallbackDetails,
    type GridRowSelectionModel,
    type GridValidRowModel,
} from '@mui/x-data-grid';
import type { FormSubmission } from './types';
import { useColumns } from './hooks/useColumns';
import { useRealtimeTableData } from './hooks/useRealtimeTableData';
import { useEffect, useMemo, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Route } from 'next';
import type { DataGridPropsWithoutDefaultValue } from '@mui/x-data-grid/internals';
import debounce from '$lib/debounce';
import { updateReadAndOrAcceptedStatus, type MyEventUpdateValueName } from './actions';
type DataGridPropsWithoutDefaultValueWithPromise<T extends GridValidRowModel> = Omit<
    DataGridPropsWithoutDefaultValue<T>,
    'onRowSelectionModelChange'
> & {
    onRowSelectionModelChange?: (
        rowSelectionModel: GridRowSelectionModel,
        details: GridCallbackDetails
    ) => Promise<void>;
};

const FormSubmissionTable = () => {
    const columns = useColumns();
    const tableData = useRealtimeTableData();
    const apiRef = useGridApiRef();
    const rows = useMemo(() => tableData, [tableData]);
    const router = useRouter();
    const [selectedRows, setSelectedRows] = useState<readonly GridRowId[]>([]);

    useEffect(() => {
        const handleRowHover: GridEventListener<'rowMouseEnter'> = (params: GridRowParams<FormSubmission>) => {
            router.prefetch(`/admin/dashboard/form-submissions/preview/${params.id}/${params.row.userId}` as Route);
        };

        return apiRef.current.subscribeEvent('rowMouseEnter', handleRowHover);
    }, [apiRef]);

    const setDebouncedSelectedRows: DataGridPropsWithoutDefaultValueWithPromise<FormSubmission>['onRowSelectionModelChange'] =
        debounce((selectedRows, details): void => {
            setSelectedRows(selectedRows);
        }, 1500);
    const handleSelectionChange: DataGridPropsWithoutDefaultValue<FormSubmission>['onRowSelectionModelChange'] = async (
        rowSelectionModel,
        details
    ) => {
        await setDebouncedSelectedRows(rowSelectionModel, details);
    };
    console.log(selectedRows);

    return (
        <DataGrid
            sx={{ insetBlockStart: '2rem', '.MuiDataGrid-cell--editable': { backgroundColor: 'secondary.dark' } }}
            rows={rows ? rows : []}
            columns={columns}
            slots={{
                toolbar: GridToolbar,
            }}
            initialState={{
                pagination: {
                    paginationModel: {
                        pageSize: 10,
                    },
                },
                sorting: {
                    sortModel: [
                        {
                            field: 'isRead',
                            sort: 'asc',
                        },
                    ],
                },
                filter: {
                    filterModel: {
                        items: [
                            {
                                field: 'isSubmitted',
                                operator: 'is',
                                value: 'true',
                            },
                        ],
                    },
                },
            }}
            processRowUpdate={async (newEdit, oldEdit) => {
                updateReadAndOrAcceptedStatus(newEdit.documentPath, {
                    isRead: newEdit.isRead,
                    isAccepted: newEdit.isAccepted,
                });
                return newEdit;
            }}
            pageSizeOptions={[10]}
            onRowSelectionModelChange={handleSelectionChange}
            loading={!!!rows}
            checkboxSelection
            disableRowSelectionOnClick
            apiRef={apiRef}
            autoHeight
        />
    );
};

export default FormSubmissionTable;
