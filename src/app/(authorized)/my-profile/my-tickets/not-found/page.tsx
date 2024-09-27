import { Paper, Typography, Box, Link } from '@mui/material';
import BuyTicketButton from '../shared/ui/BuyTicketButton';

type Props = {};

const TicketNotFound = ({ }: Props) => {
    return (
        <Box sx={{ display: 'grid', height: 'var(--centering-height)', placeContent: 'center' }}>
            <Paper sx={{ marginBottom: '2rem', paddingInline: '1rem', maxWidth: '400px' }}>
                <Typography variant="h1">Fant ingen billetter.</Typography>
                <Typography>
                    Vi fann ingen billettar registrert på denne epostadressa. Det betyr at du anten ikkje har kjøpt
                    billettar endå, eller at du har kjøpt billettane på ei anna anna epostadresse enn den du er logga
                    inn med her.
                </Typography>
                <Box sx={{ display: 'grid', gap: '1rem', marginBlock: '2rem' }}>
                    <BuyTicketButton />
                </Box>
                <Typography sx={{ marginBottom: '1rem' }}>
                    Kjøp billettar på Checkin, lag ein brukar på riktig mailadresse, eller ta kontakt med
                    <Typography component="span" sx={{ color: 'primary.main', marginInline: '1ch' }}>
                        <Link href={'mailto:styret@regncon.no'}>styret@regncon.no</Link>
                    </Typography>
                    noko er galt, eller om du ønsker billettane overført til epostadressa du har laga brukar til.
                </Typography>
            </Paper>
        </Box>
    );
};

export default TicketNotFound;
