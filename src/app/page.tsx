import Container from '@mui/material/Container';
import AppBar from '@mui/material/AppBar';
export default function Home() {
	return (
		<Container
			component={'main'}
			maxWidth="xl"
			sx={{
				display: 'grid',
				placeContent: 'center',
			}}
		></Container>
	);
}
