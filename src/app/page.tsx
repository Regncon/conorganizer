import Container from '@mui/material/Container';
import { Box } from '@mui/material';
import Navbar from './Navbar';

export default function Home() {
	return (
		<>
			<Navbar />
			<Container
				component={'main'}
				maxWidth="xl"
				sx={{
					display: 'grid',
					placeContent: 'center',
				}}
			></Container>
		</>
	);
}
