'use client';
import { useState } from 'react';
import AdultsOnlyIcon from '$lib/components/icons/AdultsOnlyIcon';
import ChildFriendlyIcon from '$lib/components/icons/ChildFriendlyIcon';
import DeleteIcon from '@mui/icons-material/Delete';
import { ActionResponse, Participant } from '$lib/types';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import {
    Accordion,
    AccordionDetails,
    AccordionSummary,
    Alert,
    Box,
    Button,
    Card,
    CardActions,
    CardContent,
    CardHeader,
    CircularProgress,
    IconButton,
    Stack,
    TextField,
    Typography,
} from '@mui/material';
import {
    ConnectEmailToParticipant,
    DeleteConnectedEmail,
} from '$app/(authorized)/my-profile/my-tickets/components/lib/actions/actions';

type Props = {
    participant: Participant;
};

const ParticipantCard = ({ participant }: Props) => {
    const [email, setEmail] = useState('');
    const [connectResponce, setConnectResponce] = useState<ActionResponse>();
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
        <Card sx={{ minWidth: 306, maxWidth: 430 }}>
            <CardHeader
                title={`${participant.firstName} ${participant.lastName}`}
                subheader={participant.ticketCategory}
            />
            <CardContent sx={{ paddingTop: 0 }}>
                <Stack direction="row" spacing={2}>
                    <Typography sx={{ fontWeight: 'bold' }}>Bestilling:</Typography>
                    <Typography> {participant.orderId}</Typography>
                </Stack>
                <Box sx={{ display: 'flex', justifyContent: 'start', alignItems: 'center' }}>
                    {participant.over18 ?
                        <>
                            <AdultsOnlyIcon chipMargin={false} />
                            <Typography sx={{ paddingLeft: '0.5rem', fontWeight: 'bold' }}>Over 18</Typography>
                        </>
                        : <>
                            <ChildFriendlyIcon chipMargin={false} />
                            <Typography sx={{ paddingLeft: '0.5rem', fontWeight: 'bold' }}>Under 18</Typography>
                        </>
                    }
                </Box>
            </CardContent>
            <CardContent>
                <Accordion sx={{ backgroundColor: '#373B57' }}>
                    <AccordionSummary expandIcon={<ExpandMoreIcon />} aria-controls="panel1-content" id="panel1-header">
                        Påmeldte arangementer
                    </AccordionSummary>
                    <AccordionDetails>Ingen påmeldte arrangementer</AccordionDetails>
                </Accordion>
                <Accordion sx={{ backgroundColor: '#373B57' }}>
                    <AccordionSummary expandIcon={<ExpandMoreIcon />} aria-controls="panel2-content" id="panel2-header">
                        Ønsker
                    </AccordionSummary>
                    <AccordionDetails>Imgen ønsker registrert</AccordionDetails>
                </Accordion>
            </CardContent>
            <CardContent>
                {loading ?
                    <CircularProgress />
                    : connectResponce && <Alert severity={connectResponce.type}>{connectResponce.message}</Alert>}
            </CardContent>
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
                    <Box sx={{ display: 'flex', justifyContent: 'start', alignItems: 'center' }}>
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
export default ParticipantCard;
