import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import Paper from '@mui/material/Paper';
import Typography from '@mui/material/Typography';
import styles from './page.module.scss';
import { Card, CardContent, CardHeader } from '@mui/material';

export default function EventCardBig() {
	return (
		<Card sx={{ backgroundImage: 'url(/blekksprut.jpg)' }}>
			<CardHeader title="hei1" />
			<CardContent>hei</CardContent>
		</Card>
	);
}
