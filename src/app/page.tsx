import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import Paper from '@mui/material/Paper';
import Typography from '@mui/material/Typography';
import styles from './page.module.scss';
import EventCardBig from './EventCardBig';

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
			<Box
				component={Paper}
				sx={{ display: 'grid', placeContent: 'center' }}
				className={styles['main-test']}
				elevation={1}
			>
			<img src="/placeholderlogo.png" alt="logo" />
			<EventCardBig/>
			</Box>
		</Container>
	);
}
