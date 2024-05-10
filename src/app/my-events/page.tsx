import CardBase from '$app/dashboard/CardBase';
import Paper from '@mui/material/Paper';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';

const MyEvents = async () => {
	return (
		<Paper>
			<Grid2 container justifyContent="space-between">
				<Grid2 xs={5.7}>
					<CardBase
						href="/my-events"
						description="Trykk for og gå til mine arrangementer"
						img="/my-events.jpg"
						imgAlt="Mine arrangementer"
						title="Mine arrangementer"
					/>
				</Grid2>
				<Grid2 xs={5.7}>
					<CardBase
						href="/my-events"
						description="Trykk for og gå til mine arrangementer"
						img="/my-events.jpg"
						imgAlt="Mine arrangementer"
						title="Mine arrangementer"
					/>
				</Grid2>
			</Grid2>
		</Paper>
	);
};

export default MyEvents;
