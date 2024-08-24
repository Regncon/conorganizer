'use client';
import { DataGrid, useGridApiRef, type GridEventListener, type GridRowParams, GridToolbar } from '@mui/x-data-grid';
import type { FormSubmission } from './types';
import { useColumns } from './hooks/useColumns';
import { useRealtimeTableData } from './hooks/useRealtimeTableData';
import { useEffect, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { Route } from 'next';

const FormSubmissionTable = () => {
    const columns = useColumns();
    const tableData = useRealtimeTableData();
    const apiRef = useGridApiRef();
    const rows = useMemo(() => tableData, [tableData]);
    const router = useRouter();

    useEffect(() => {
        const handleRowHover: GridEventListener<'rowMouseEnter'> = (params: GridRowParams<FormSubmission>) => {
            router.prefetch(`/admin/dashboard/form-submissions/preview/${params.id}` as Route);
        };

        return apiRef.current.subscribeEvent('rowMouseEnter', handleRowHover);
    }, [apiRef]);

    return (
        <DataGrid
            sx={{ insetBlockStart: '2rem' }}
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
            pageSizeOptions={[10]}
            onPaginationModelChange={(e) => {
                console.log(e);
            }}
            loading={!!!rows}
            // checkboxSelection
            disableRowSelectionOnClick
            apiRef={apiRef}
            autoHeight
        />
    );
};

export default FormSubmissionTable;
