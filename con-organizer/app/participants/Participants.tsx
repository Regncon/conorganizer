import { MouseEvent, useEffect, useState } from 'react';
import { LoadingButton } from '@mui/lab';
import Box from '@mui/material/Box';
import { DataGrid, GridColDef } from '@mui/x-data-grid';
import { useAuth } from '@/components/AuthProvider';
import { useAllParticipants } from '@/lib/hooks/UseAllParticipants';
import { useAllUserSettings } from '@/lib/hooks/UseAllUserSettings';
import { useUserSettings } from '@/lib/hooks/UseUserSettings';
import { Participant } from '@/models/types';
import { Pool } from '@/models/enums';

type NewParticipants = {
    newParticipants: Participant[];
};

const Participants = () => {
    const user = useAuth();
    const { allUserSettings, loadingUserSettings } = useAllUserSettings();

    const { participants, loadingParticipants } = useAllParticipants();

    const [rows, setRows] = useState<Participant[]>([]);

    const { userSettings } = useUserSettings(user?.uid);
    const [isAdmin, setIsAdmin] = useState<boolean>(false);

    useEffect(() => {
        setIsAdmin(userSettings?.admin && user ? true : false);
    }, [user, userSettings]);

    const columns: GridColDef[] = [
        /*         { field: 'id', headerName: 'ID', width: 90 }, */
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
    ];

    useEffect(() => {
        if (participants) {
            setRows(participants);
        }
    }, [participants, setRows]);

    const [newParticipants, setNewParticipants] = useState<Participant[]>([]);

    const [loadingParticipantsFromCheckIn, setLoadingParticipantsFromCheckIn] = useState<boolean>(false);
    const getParticipantsFromCheckIn = async (e: MouseEvent<HTMLButtonElement>) => {
        e.preventDefault();
        setLoadingParticipantsFromCheckIn(true);

        try {
            const result = await fetch('/api/getParticipants', { method: 'GET' });
            const res = (await result.json()) as NewParticipants;
            console.log(res.newParticipants, 'res newParticipants');
            setNewParticipants(res.newParticipants);
        } catch (error) {
            console.log(error);
        }

        setLoadingParticipantsFromCheckIn(false);
    };

    const [migrateParticipants, setMigrateParticipants] = useState<boolean>(false);

    const [migratedParticipantsList, setMigratedParticipantsList] = useState<Participant[]>([]);

    const callMigrateParticipants = async (e: MouseEvent<HTMLButtonElement>) => {
        e.preventDefault();
        setMigrateParticipants(true);

        try {
            const result = await fetch('/api/migrateParticipants', { method: 'GET' });
            const res = (await result.json()) as NewParticipants;
            console.log(res.newParticipants, 'res newParticipants');
            setMigratedParticipantsList(res.newParticipants);
        } catch (error) {
            console.log(error);
        }

        setMigrateParticipants(false);
    };

    const [loadingGenerateEnrollmentList, setLoadingGenerateEnrollmentList] = useState<boolean>(false);
    const generateEnrollmentList = async (e: MouseEvent<HTMLButtonElement>, pool: Pool ) => {
        e.preventDefault();
        setLoadingGenerateEnrollmentList(true);

        try {
            const result = await fetch('/api/generateEnrollmentList', { method: 'POST', body: JSON.stringify({ pool }) });
            //const res = (await result.json());
        } catch (error) {
            console.log(error);
        }

        setLoadingGenerateEnrollmentList(false);
    };


    useEffect(() => {
        console.log(newParticipants, 'newParticipants array');
    }, [newParticipants]);

    return (
        <Box sx={isAdmin ? { maxWidth: 1000, margin: 'auto' } : { display: 'none' }}>
            <h1>Deltakere</h1>

            <Box
                sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    marginBottom: '1em',
                }}
            >
                <LoadingButton
                    loading={loadingParticipantsFromCheckIn}
                    variant="contained"
                    onClick={(e) => getParticipantsFromCheckIn(e)}
                >
                    Hent deltakere fra CheckIn
                </LoadingButton>
                <LoadingButton
                    loading={migrateParticipants}
                    variant="contained"
                    onClick={(e) => callMigrateParticipants(e)}
                >
                    Migrer deltakere til nytt format
                </LoadingButton>
            </Box>

            <Box
                sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    marginBottom: '1em',
                }}
            >
                <LoadingButton
                    loading={loadingGenerateEnrollmentList}
                    variant="outlined"
                    onClick={(e) => generateEnrollmentList(e, Pool.FridayEvening)}
                >
                    Hent Puljedata for fredag kveld
                </LoadingButton>
            </Box>
            <Box sx={newParticipants?.length > 0 ? { display: 'block' } : { display: 'none' }}>
                <h2>{newParticipants?.length} Nye deltagere lagt til i basen</h2>
                {newParticipants?.map((participant) => (
                    <p key={participant.externalId}>{participant.name}</p>
                ))}
            </Box>

            <Box sx={migratedParticipantsList?.length > 0 ? { display: 'block' } : { display: 'none' }}>
                <h2>{migratedParticipantsList?.length} Deltagere migrert til nytt format</h2>
                {migratedParticipantsList?.map((participant) => (
                    <p key={participant.externalId}>{participant.name}</p>
                ))}
            </Box>

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
