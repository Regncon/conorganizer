import { Paper, Typography } from '@mui/material';
import Grid2 from '@mui/material/Unstable_Grid2';
import EventForm from './EventForm';
type Props = {
    params: {
        id: string;
    };
    searchParams: {
        userId: string;
    };
};
const Create = async ({ params: { id } }: Props) => {
    return (
        <>
            <Paper sx={{ padding: '1rem', marginBottom: '2rem' }}>
                <Grid2 container gap="3rem">
                    <Typography variant="h1">Meld på arrangement til Regncon XXXII 2024</Typography>
                    <Typography>
                        Takk for at du vil arrangere eit spel på Regncon, anten det er brettspel, kortspel, rollespel
                        eller anna, så sett vi enormt pris på ditt bidrag. Fyll inn skjemaet så godt du kan, og ikkje
                        vere redd for å ta kontakt med Regnconstyret på{' '}
                        <a href="mailto:regncon@gmail.com">regncon@gmail.com</a> om du skulle lure på noko!
                    </Typography>
                </Grid2>
            </Paper>
            <EventForm id={id} />
        </>
    );
};
export default Create;
