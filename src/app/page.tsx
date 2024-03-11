import Container from '@mui/material/Container';

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
