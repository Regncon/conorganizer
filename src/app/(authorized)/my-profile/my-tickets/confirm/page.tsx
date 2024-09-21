import { Paper, Typography, Box } from '@mui/material';
import ConfirmEmailButtons from './components/ui/ConfirmEmailButtons';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import BuyTicketButton from '../shared/ui/BuyTicketButton';

type Props = {};

const Confirm = async ({ }: Props) => {
    const { user } = await getAuthorizedAuth();

    return (
        <Box sx={{ display: 'grid', placeContent: 'center', height: 'var(--centering-height)' }}>
            <Paper sx={{ marginBottom: '2rem', paddingInline: '0.5rem', maxWidth: '400px', padding: '2rem' }}>
                <Typography variant="h1">Bekreft e-post</Typography>
                <Typography sx={{ margin: '1rem' }}>
                    Du må bekrefte e-posten din for å få tilgang til dine billetter. Trykk på knappen under for å sende
                    en bekreftelses e-post. Sjekk spam-mappen hvis du ikke finner e-posten.
                </Typography>
                <Box sx={{ display: 'grid', gap: '1rem' }}>
                    <ConfirmEmailButtons disabled={user?.emailVerified} />
                </Box>
                <Typography variant="h2">Har ikke billetter</Typography>
                <BuyTicketButton />
            </Paper>
        </Box>
    );
};

export default Confirm;
