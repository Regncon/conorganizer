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
    GridPagination,
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
import { Box, Button, CircularProgress } from '@mui/material';
import { convertToConEvent } from './preview/[id]/[userid]/actions';
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
    const [selectedRowsModel, setSelectedRowsModel] = useState<readonly GridRowId[]>([]);
    const [isLoadingConverting, setIsLoadingConverting] = useState<boolean>(false);

    useEffect(() => {
        const handleRowHover: GridEventListener<'rowMouseEnter'> = (params: GridRowParams<FormSubmission>) => {
            router.prefetch(`/admin/dashboard/form-submissions/preview/${params.id}/${params.row.userId}` as Route);
        };

        return apiRef.current.subscribeEvent('rowMouseEnter', handleRowHover);
    }, [apiRef]);

    const debouncedSetSelectedRows: DataGridPropsWithoutDefaultValueWithPromise<FormSubmission>['onRowSelectionModelChange'] =
        debounce((rowSelectionModel, details): void => {
            setSelectedRowsModel(rowSelectionModel);
        }, 800);

    const handleSelectionChange: DataGridPropsWithoutDefaultValueWithPromise<FormSubmission>['onRowSelectionModelChange'] =
        async (rowSelectionModel, details) => {
            await debouncedSetSelectedRows(rowSelectionModel, details);
            details.api.getRow;
        };

    const handleConvertToEvents = async () => {
        setIsLoadingConverting(true);
        const selectedRows = selectedRowsModel.map((rowId) => {
            const row: FormSubmission | null = apiRef.current.getRow(rowId);
            if (row?.id === undefined || row?.userId === undefined) {
                return null;
            }
            return () => convertToConEvent(row.id, row.userId);
        });

        await Promise.all(selectedRows.map((row) => row?.()));
        setIsLoadingConverting(false);
    };
    console.log(selectedRowsModel);
    return (
        <DataGrid
            sx={{ insetBlockStart: '2rem', '.MuiDataGrid-cell--editable': { backgroundColor: 'secondary.dark' } }}
            rows={rows ? rows : []}
            columns={columns}
            slots={{
                toolbar: GridToolbar,
                footer: () => (
                    <>
                        <GridPagination />
                        {selectedRowsModel.length > 0 ?
                            <Box sx={{ display: 'grid', placeContent: 'end' }}>
                                <Button
                                    variant="contained"
                                    disabled={isLoadingConverting}
                                    onClick={handleConvertToEvents}
                                >
                                    {isLoadingConverting ?
                                        <>
                                            Vennligst vent Konverter arrangementer
                                            <CircularProgress sx={{ marginInlineStart: '2rem' }} size="1.5rem" />
                                        </>
                                    :   'Konverter til arrangementer'}
                                </Button>
                            </Box>
                        :   null}
                    </>
                ),
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
