import Container from '@mui/material/Container';
import Landing from './Landing';

export default function Home() {
	return (
		<>
			<Landing />
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
