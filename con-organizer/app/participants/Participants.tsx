import * as React from 'react';
import Box from '@mui/material/Box';
import { DataGrid, GridColDef } from '@mui/x-data-grid';
import { useAllParticipants } from '@/lib/hooks/UseAllParticipants';
import { useAllUserSettings } from '@/lib/hooks/UseAllUserSettings';
import { Participant } from '@/models/types';

const Participants = () => {

    const {userSettings, loadingUserSettings} = useAllUserSettings();

    console.log(userSettings);

    const {participants, loadingParticipants} = useAllParticipants();
    console.log(participants);


    const columns: GridColDef[] = [
        { field: 'id', headerName: 'ID', width: 90 },
        {
            field: 'name',
            headerName: 'Navn',
            width: 150,
            editable: true, 
        },
        {
            field: 'email',
            headerName: 'Epost',
            width: 150,
            editable: false,
        },
        {
            field: 'isChild',
            headerName: 'Barn',
            type: 'boolean',
            width: 110,
            editable: false,
        },
        {
            field: 'type',
            headerName: 'Type',
            width: 110,
            editable: false,
        },
        {
            field: 'connectedUser',
            headerName: 'Koblet til',
            width: 110,
            editable: false,
        },
        {
            field: 'ticketType',
            headerName: 'Billettype',
            width: 110,
            editable: false,
        },
    ];

    const rows = [
        { id: 1, name: 'Ola Norman', email: 'ola@test.com', isChild: false, type: 'Bruker', connectedUser: '' },
        { id: 2, name: 'Kari Norman', email: 'kari@test.com', isChild: false, type: 'Bruker', connectedUser: '' },
        { id: 3, name: 'Truls Norman', email: 'ola@test.com', isChild: true, type: 'Deltaker', connectedUser: 'Ola Norman' },
    ];

    return (
        <Box sx={{ maxWidth: 1000, margin: 'auto' }}>
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
                pageSizeOptions={[5]}
                disableRowSelectionOnClick
                
            />
        </Box>
    );
};

export default Participants;
