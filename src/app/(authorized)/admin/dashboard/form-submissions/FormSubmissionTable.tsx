'use client';
import { DataGrid, useGridApiRef, type GridComparatorFn, type GridRowParams } from '@mui/x-data-grid';
import type { FormSubmission, FormSubmissionColDef } from './types';
import { useRealtimeTableData } from './hooks';
import { useMemo, useRef } from 'react';
import { updateReadStatus } from './actions';

const columns: FormSubmissionColDef = [
    { field: 'id', headerName: 'ID', width: 255 },
    {
        field: 'title',
        headerName: 'Tittel',
        width: 140,
        editable: true,
    },
    {
        field: 'subTitle',
        headerName: 'Kort beskrivelse',
        width: 240,
        editable: true,
    },
    {
        field: 'isRead',
        headerName: 'Lest',
        type: 'boolean',
        width: 110,
        editable: true,
    },
    {
        field: 'isAccepted',
        headerName: 'Godkjent',
        type: 'boolean',
        width: 110,
        editable: true,
    },
];

// const rows: FormSubmission[] = [
//     {
//         id: '377104b0-7354-5125-81c1-167ddce1df1e',
//         title: 'further',
//         subTitle: 'stay elephant wore rocky him cattle describe shelf itself activity ',
//         isRead: true,
//         isAccepted: false,
//     },
//     {
//         id: '2dc4d222-a50c-5ddd-b607-7df06abcd36e',
//         title: 'saw',
//         subTitle: 'horse imagine saw courage way call ask brush has tongue ',
//         isRead: true,
//         isAccepted: false,
//     },
//     {
//         id: 'd6bf7fb7-9b0e-5075-aa94-5cc1c312ec31',
//         title: 'nearest',
//         subTitle: 'dish observe total shells leg sharp elephant put whole report ',
//         isRead: true,
//         isAccepted: false,
//     },
//     {
//         id: '5932bf9a-ddff-5b9a-99dc-561240fbb8b8',
//         title: 'way',
//         subTitle: 'score beautiful roar danger nine expression far twenty rise stronger ',
//         isRead: true,
//         isAccepted: false,
//     },
//     {
//         id: '39531de4-08c5-5ccb-99eb-a45228e3f063',
//         title: 'climb',
//         subTitle: 'reach hurried successful dirty egg fall victory cell track although ',
//         isRead: true,
//         isAccepted: false,
//     },
//     {
//         id: '191d972a-20aa-52b5-a359-1546ec517973',
//         title: 'title',
//         subTitle: 'shelf very generally invented pie promised type gift is outside ',
//         isRead: true,
//         isAccepted: false,
//     },
//     {
//         id: '3bf5a011-c2eb-5e45-9193-896cb4060af9',
//         title: 'slope',
//         subTitle: 'arrow feed facing putting many telephone ants pale wire more ',
//         isRead: true,
//         isAccepted: false,
//     },
//     {
//         id: '6f0c5d97-e683-5c21-b9f2-123b70ab4c18',
//         title: 'note',
//         subTitle: 'below glass engine fence sell were town passage sister spent ',
//         isRead: true,
//         isAccepted: false,
//     },
//     {
//         id: '67dfe3c9-62d0-5378-aba0-ddc6a33bc46f',
//         title: 'same',
//         subTitle: 'kids took powerful national live football memory has select doing ',
//         isRead: true,
//         isAccepted: false,
//     },
// ];

const FormSubmissionTable = () => {
    const tableData = useRealtimeTableData();
    const apiRef = useGridApiRef();
    const rows = useMemo(() => tableData, [tableData]);

    return (
        <DataGrid
            sx={{ minHeight: '40rem', insetBlockStart: '2rem' }}
            rows={rows}
            columns={columns}
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
            }}
            // onRowSelectionModelChange={(e, b) => {
            //     console.log(e, 'e', b.api.getRow(b.api.getAllRowIds()[0]), 'b');
            // }}
            onRowClick={(e: GridRowParams<FormSubmission>) => {
                console.log(e.row);
                updateReadStatus(e.row.id);
            }}
            pageSizeOptions={[10]}
            onPaginationModelChange={(e) => {
                console.log(e);
            }}
            loading={rows.length === 0}
            checkboxSelection
            disableRowSelectionOnClick
            apiRef={apiRef}
        />
    );
};

export default FormSubmissionTable;
