import FormSubmissions from './components/FormSubmissions';
import Events from './components/Events';
import MyEvents from '$app/(authorized)/dashboard/components/MyEvents';
import { Box } from '@mui/material';
import CardBase from '$app/(authorized)/dashboard/components/CardBase';

const Dashboard = async () => {
    return (
        <Box
            sx={{
                display: 'grid',
                gridTemplateColumns: 'repeat(auto-fit, minmax(16.7rem, 0.2fr))',
                gap: 2,
                placeItems: 'center',
                placeContent: 'center',
                marginBlockStart: '1rem',
            }}
        >
            <FormSubmissions />
            <CardBase
                href="/admin/dashboard/rooms?pool=fridayEvening"
                subTitle="Trykk for å gå til romfordelingen"
                img="/rooms-small.webp"
                imgAlt="Romfordeling"
                title="Romfordeling"
            />
            <Events />
            <MyEvents />
        </Box>
    );
};

export default Dashboard;
