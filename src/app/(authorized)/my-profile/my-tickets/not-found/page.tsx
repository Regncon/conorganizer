import { Paper, Typography, Box } from '@mui/material';
import BuyTicketButton from '../shared/ui/BuyTicketButton';

type Props = {};

const TicketNotFound = ({ }: Props) => {
    return (
        <Box sx={{ display: 'grid', height: 'var(--centering-height)', placeContent: 'center' }}>
            <Paper sx={{ marginBottom: '2rem', paddingInline: '1rem', maxWidth: '400px' }}>
                <Typography variant="h1">Fant ingen billetter.</Typography>
                <Typography>
                    Legge inn en fin hjelpetekst skrevet av en som ikke er meg eller dyslektiker som forklarer at folk
                    må ta kontakt dersom de trenger hjelp
                </Typography>
                <Box sx={{ display: 'grid', gap: '1rem', marginBlockEnd: '1rem' }}>
                    <BuyTicketButton />
                </Box>
            </Paper>
        </Box>
    );
};

export default TicketNotFound;
