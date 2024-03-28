import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import Paper from '@mui/material/Paper';
import Typography from '@mui/material/Typography';
import styles from './page.module.scss';
import EventCardBig from './EventCardBig';
import EventCardSmall from './EventCardSmall';

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
				<EventCardBig title='Hello world' gameMaster='Gerhard Fajita' shortDescription='Mord overalt! Kos! Gøy!' system='Call of Chthuhlth'/>
				<Box sx={{ display: 'flex' }}>
					<EventCardSmall title='Hi' gameMaster='Gardh Fajita2' system='Dungeons 2'/>
					<EventCardSmall title='Any% speedrun' gameMaster='Gorde Fajita3' system='Terraria'/>
				</Box>
			</Box>
		</Container>
	);
}
