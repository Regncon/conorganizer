import { Paper, Typography, Button, TextField, Box, FormControl } from '@mui/material';
import ConfirmEmailButton from './UI/ConfirmEmail';

type Props = {};

const TicketNotFound = ({}: Props) => {
    return (
        <Paper sx={{ marginBottom: '2rem', paddingInline: '1rem', maxWidth: '320px' }}>
            <Typography variant="h1">Ingen?/Mine billetter</Typography>
            <Typography variant="h2">Fant ingen billetter.</Typography>
            <Box sx={{ display: 'grid', gap: '1rem', marginBlockEnd: '1rem' }}>
                <Button fullWidth variant="contained">
                    Kjøp billett
                </Button>
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
    );
};

export default TicketNotFound;
