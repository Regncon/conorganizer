import AdultsOnlyIcon from '$lib/components/icons/AdultsOnlyIcon';
import ChildFriendlyIcon from '$lib/components/icons/ChildFriendlyIcon';
import { Participant } from '$lib/types';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import {
    Accordion,
    AccordionActions,
    AccordionDetails,
    AccordionSummary,
    Box,
    Button,
    Card,
    CardActions,
    CardContent,
    CardHeader,
    Stack,
    Typography,
} from '@mui/material';

type Props = {
    participant: Participant;
};

const ParticipantCard = ({ participant }: Props) => {
    return (
        <Card sx={{ minWidth: 306 }}>
            <CardHeader title={participant.name} subheader={participant.ticketType} />
            <CardContent sx={{ paddingTop: 0 }}>
                <Box sx={{ display: 'flex', justifyContent: 'start', alignItems: 'center' }}>
                    {participant.over18 ?
                        <>
                            <Typography sx={{ fontWeight: 'bold' }}>Over 18</Typography>
                            <AdultsOnlyIcon />
                        </>
                        : <>
                            <Typography sx={{ fontWeight: 'bold' }}>Under 18</Typography>
                            <ChildFriendlyIcon />
                        </>
                    }
                </Box>
                <Stack direction="row" spacing={2}>
                    <Typography sx={{ fontWeight: 'bold' }}>Bestilling:</Typography>
                    <Typography> {participant.ticketId}</Typography>
                    <Typography sx={{ fontWeight: 'bold' }}>Status:</Typography>
                    <Typography> {participant.ticketStatus}</Typography>
                </Stack>
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
                <Accordion sx={{ backgroundColor: '#373B57' }}>
                    <AccordionSummary expandIcon={<ExpandMoreIcon />} aria-controls="panel3-content" id="panel3-header">
                        Notater
                    </AccordionSummary>
                    <AccordionDetails>
                        Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse malesuada lacus ex, sit
                        amet blandit leo lobortis eget.
                    </AccordionDetails>
                    <AccordionActions>
                        <Button>Cancel</Button>
                        <Button>Lagre</Button>
                    </AccordionActions>
                </Accordion>
            </CardContent>
            <CardActions>
                <Button size="small">Learn More</Button>
            </CardActions>
        </Card>
    );
};
export default ParticipantCard;
