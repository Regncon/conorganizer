import { redirectToAdminDashboardWhenAdministrator } from '$lib/libServer';
import MyEvents from './MyEvents';

import { Grid2 } from '@mui/material';

const Dashboard = async () => {
    await redirectToAdminDashboardWhenAdministrator();
    return (
        <Grid2 container spacing="2rem" sx={{ marginTop: '0.5rem' }}>
            {/* <Grid2 xs={12} md={3}>
                <MyTickets />
            </Grid2> */}
            <Grid2
                size={{
                    xs: 12,
                    md: 3,
                }}
            >
                <MyEvents />
            </Grid2>
        </Grid2>
    );
};

export default Dashboard;
