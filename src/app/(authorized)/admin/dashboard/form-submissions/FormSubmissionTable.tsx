'use client';
import { DataGrid } from '@mui/x-data-grid';
import type { FromSubmission, FromSubmissionColDef } from './types';

const columns: FromSubmissionColDef = [
    { field: 'id', headerName: 'ID', width: 255 },
    {
        field: 'title',
        headerName: 'Tittel',
        width: 140,
        editable: true,
    },
    {
        field: 'subtitle',
        headerName: 'Kort beskrivelse',
        width: 250,
        editable: true,
    },
    {
        field: 'isRead',
        headerName: 'Lest',
        type: 'number',
        width: 110,
        editable: true,
    },
    {
        field: 'isAccepted',
        headerName: 'Godkjent',
        type: 'number',
        width: 110,
        editable: true,
    },
];

const rows: FromSubmission[] = [
    {
        id: '377104b0-7354-5125-81c1-167ddce1df1e',
        title: 'further',
        subtitle: 'stay elephant wore rocky him cattle describe shelf itself activity ',
        isRead: 'true',
        isAccepted: 'false',
    },
    {
        id: '2dc4d222-a50c-5ddd-b607-7df06abcd36e',
        title: 'saw',
        subtitle: 'horse imagine saw courage way call ask brush has tongue ',
        isRead: 'true',
        isAccepted: 'false',
    },
    {
        id: 'd6bf7fb7-9b0e-5075-aa94-5cc1c312ec31',
        title: 'nearest',
        subtitle: 'dish observe total shells leg sharp elephant put whole report ',
        isRead: 'true',
        isAccepted: 'false',
    },
    {
        id: '5932bf9a-ddff-5b9a-99dc-561240fbb8b8',
        title: 'way',
        subtitle: 'score beautiful roar danger nine expression far twenty rise stronger ',
        isRead: 'true',
        isAccepted: 'false',
    },
    {
        id: '39531de4-08c5-5ccb-99eb-a45228e3f063',
        title: 'climb',
        subtitle: 'reach hurried successful dirty egg fall victory cell track although ',
        isRead: 'true',
        isAccepted: 'false',
    },
    {
        id: '191d972a-20aa-52b5-a359-1546ec517973',
        title: 'title',
        subtitle: 'shelf very generally invented pie promised type gift is outside ',
        isRead: 'true',
        isAccepted: 'false',
    },
    {
        id: '3bf5a011-c2eb-5e45-9193-896cb4060af9',
        title: 'slope',
        subtitle: 'arrow feed facing putting many telephone ants pale wire more ',
        isRead: 'true',
        isAccepted: 'false',
    },
    {
        id: '6f0c5d97-e683-5c21-b9f2-123b70ab4c18',
        title: 'note',
        subtitle: 'below glass engine fence sell were town passage sister spent ',
        isRead: 'true',
        isAccepted: 'false',
    },
    {
        id: '67dfe3c9-62d0-5378-aba0-ddc6a33bc46f',
        title: 'same',
        subtitle: 'kids took powerful national live football memory has select doing ',
        isRead: 'true',
        isAccepted: 'false',
    },
];

const FormSubmissionTable = () => {
    return (
        <DataGrid
            rows={rows}
            columns={columns}
            initialState={{
                pagination: {
                    paginationModel: {
                        pageSize: 5,
                    },
                },
            }}
            // onRowSelectionModelChange={(e, b) => {
            //     console.log(e, 'e', b.api.getRow(b.api.getAllRowIds()[0]), 'b');
            // }}
            pageSizeOptions={[5]}
            checkboxSelection
            disableRowSelectionOnClick
        />
    );
};

export default FormSubmissionTable;
