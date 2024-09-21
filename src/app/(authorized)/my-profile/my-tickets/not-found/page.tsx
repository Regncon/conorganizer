import { Paper, Typography, Box, Button, TextField } from '@mui/material';
import BuyTicketButton from '../shared/ui/BuyTicketButton';

type Props = {};

const TicketNotFound = ({}: Props) => {
    return (
        <Box sx={{ display: 'grid', height: 'var(--centering-height)', placeContent: 'center' }}>
            <Paper sx={{ marginBottom: '2rem', paddingInline: '1rem', maxWidth: '320px' }}>
                <Typography variant="h1">Ingen?/Mine billetter</Typography>
                <Typography variant="h2">Fant ingen billetter.</Typography>
                <Typography>
                    Legge inn en fin hjelpetekst skrevet av en som ikke er meg eller dyslektiker som forklarer at folk
                    må ta kontakt dersom de trenger hjelp
                </Typography>
                <Box sx={{ display: 'grid', gap: '1rem', marginBlockEnd: '1rem' }}>
                    <BuyTicketButton />
                    <Button fullWidth variant="contained">
                        Har allerede kjøpt billett
                    </Button>
                </Box>

                <Box sx={{ display: 'grid', gap: '1rem', marginBlockEnd: '1rem' }}>
                    <TextField fullWidth label="Skriv inn e-posten du brukte på Checkin" />
                    <Button fullWidth variant="contained">
                        Hent billett
                    </Button>
                </Box>
            </Paper>
        </Box>
    );
};

export default TicketNotFound;
