import Paper from '@mui/material/Paper';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import FormSubmissions from './FormSubmissions';
import Events from './Events';

const Dashboard = () => {
    return (
        <Grid2 container spacing="2rem" sx={{ marginTop: '0.5rem' }}>
            {/* <Grid2 xs={12} md={3}>
                <MyTickets />
            </Grid2> */}
            <Grid2 xs={12} md={3}>
                <FormSubmissions />
            </Grid2>
            <Grid2 xs={12} md={3}>
                <Events />
            </Grid2>
        </Grid2>
    );
};

export default Dashboard;
