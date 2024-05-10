import Paper from '@mui/material/Paper';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';
import MyTickets from './MyTickets';
import MyEvents from './MyEvents';

const Dashboard = () => {
	return (
		<Paper sx={{ marginBlock: 2 }}>
			<Grid2 container justifyContent="space-between">
				<Grid2 xs={5.7}>
					<MyTickets />
				</Grid2>
				<Grid2 xs={5.7}>
					<MyEvents />
				</Grid2>
			</Grid2>
		</Paper>
	);
};

export default Dashboard;
