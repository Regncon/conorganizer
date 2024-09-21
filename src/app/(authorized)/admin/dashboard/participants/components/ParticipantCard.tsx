'use client';
import { useState } from 'react';
import AdultsOnlyIcon from '$lib/components/icons/AdultsOnlyIcon';
import ChildFriendlyIcon from '$lib/components/icons/ChildFriendlyIcon';
import { Participant } from '$lib/types';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import {
    Accordion,
    AccordionDetails,
    AccordionSummary,
    Box,
    Button,
    Card,
    CardActions,
    CardContent,
    CardHeader,
    Stack,
    TextField,
    Typography,
} from '@mui/material';

type Props = {
    participant: Participant;
};

const ParticipantCard = ({ participant }: Props) => {
    const [email, setEmail] = useState('');

    const handleConnectToEmail = () => {
        console.log('Connect to email:', email);
    };

    return (
        <Card sx={{ minWidth: 306 }}>
            <CardHeader title={participant.name} subheader={participant.ticketCategory} />
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
            <CardActions>
                <TextField label="Epost" type="email" value={email} onChange={(e) => setEmail(e.target.value)} />
                <Button variant="contained" onClick={handleConnectToEmail}>
                    Koble til epost
                </Button>
            </CardActions>
            <CardContent>
                <Typography>Bilett epost: {participant.ticketEmail}</Typography>
                {participant.oredrEmails?.map((email, index) => (
                    <Typography key={index}>
                        Bestillings epost {index + 1}: {email}
                    </Typography>
                ))}
            </CardContent>
        </Card>
    );
};
export default ParticipantCard;
