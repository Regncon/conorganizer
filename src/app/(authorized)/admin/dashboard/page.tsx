import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import FormSubmissions from './FormSubmissions';
import Events from './Events';
import MyEvents from '$app/(authorized)/dashboard/MyEvents';

const Dashboard = async () => {
    return (
        <Grid2 container spacing="2rem" sx={{ marginTop: '0.5rem' }}>
            <Grid2 xs={12} md={3}>
                <FormSubmissions />
            </Grid2>
            {/*
            <Grid2 xs={12} md={3}>
                <Events />
            </Grid2>
            */}
            <Grid2 xs={12} md={3}>
                <MyEvents />
            </Grid2>
        </Grid2>
    );
};

export default Dashboard;
