'use client';
import { ActionResponse, Participant } from '$lib/types';
import {
    Alert,
    Box,
    Button,
    Card,
    CardActions,
    CardContent,
    CircularProgress,
    IconButton,
    TextField,
    Typography,
} from '@mui/material';
import { useState } from 'react';
import { ConnectEmailToParticipant, DeleteConnectedEmail } from '../lib/actions/actions';
import DeleteIcon from '@mui/icons-material/Delete';

type Props = {
    participant: Participant;
};

const MyParticipant = ({ participant }: Props) => {
    const [email, setEmail] = useState('');
    const [connectResponse, setConnectResponce] = useState<ActionResponse>();
    const [loading, setLoading] = useState(false);

    const handleConnectToEmail = async () => {
        setLoading(true);
        const responce = await ConnectEmailToParticipant(participant.id as string, email);
        setConnectResponce(responce);
        setLoading(false);
    };

    const handleDeleteConnected = async (email: string) => {
        setLoading(true);
        const responce = await DeleteConnectedEmail(participant.id as string, email);
        setConnectResponce(responce);
        setLoading(false);
    };

    return (
        <Card>
            <CardContent>
                <Typography>{participant.ticketCategory}</Typography>
                <Typography>{`${participant.firstName} ${participant.lastName}`}</Typography>
                <Typography>Bestilling: {participant.orderId}</Typography>
            </CardContent>

            {loading ?
                <CircularProgress />
                : connectResponse && <Alert severity={connectResponse.type}>{connectResponse.message}</Alert>}
            <CardActions>
                <TextField label="Epost" type="email" value={email} onChange={(e) => setEmail(e.target.value)} />
                <Button variant="contained" onClick={handleConnectToEmail}>
                    Koble til epost
                </Button>
            </CardActions>
            <CardContent>
                <Typography>Bilett epost: {participant.ticketEmail}</Typography>
                {participant.orderEmails?.map((email, index) => (
                    <Typography key={index}>
                        Bestillings epost {index + 1}: {email}
                    </Typography>
                ))}
                {participant.connectedEmails?.map((email, index) => (
                    <Box sx={{ display: 'flex', justifyContent: 'start', alignItems: 'center' }} key={index}>
                        <Typography>
                            Kobledt epost {index + 1}: {email}
                        </Typography>
                        <IconButton aria-label="delete" color="primary" onClick={() => handleDeleteConnected(email)}>
                            <DeleteIcon />
                        </IconButton>
                    </Box>
                ))}
            </CardContent>
        </Card>
    );
};
export default MyParticipant;
