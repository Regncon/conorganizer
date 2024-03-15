import Container from '@mui/material/Container';
import AppBar from '@mui/material/AppBar';
import { Box } from '@mui/material';
import navbar from './navbar';

export default function Home() {
	return (
		<Container
			component={'main'}
			maxWidth="xl"
			sx={{
				display: 'grid',
				placeContent: 'center',
			}}
		>
			<navbar />
		</Container>
	);
}
