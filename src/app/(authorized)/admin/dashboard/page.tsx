import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import FormSubmissions from './FormSubmissions';
import Events from './Events';
import MyEvents from '$app/(authorized)/dashboard/MyEvents';
import { Box } from '@mui/material';

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
            <Events />
            <MyEvents />
        </Box>
    );
};

export default Dashboard;
