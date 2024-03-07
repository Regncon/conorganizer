import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import Paper from '@mui/material/Paper';
import Typography from '@mui/material/Typography';
import styles from './page.module.scss';
import { Card, CardContent, CardHeader } from '@mui/material';

export default function EventCardBig() {
	return (
		<Card
			sx={{ backgroundImage: 'url(/blekksprut.jpg)', height: '267px', width: '306px', backgroundSize: 'cover' }}
		>
			<CardHeader title="Tentakkel" sx={{ height: '141px',alignItems:"flex-end" }} />
			<CardContent
				sx={{ height: '126px', backgroundColor: 'rgba(0,0,0,0.5)', backdropFilter: 'blur(4px)', padding: '0' }}
			>
				<Typography> Gerhard Fajita </Typography>
				<Box>
					<Typography> Call of Cthulhu </Typography>
				</Box>
				<Typography sx={{ color: 'white' }}>
					Klarer du å overleve landets farligste fjell? Pass på at du ikke brenner deg!
				</Typography>
			</CardContent>
		</Card>
	);
}
