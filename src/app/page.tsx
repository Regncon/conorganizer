import Container from '@mui/material/Container';
import styles from './page.module.scss';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';

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
				height={'20rem'}
				width={'50rem'}
			>
				<Typography>Hello world! :D</Typography>
			</Box>
		</Container>
	);
}
