import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import { Card, CardContent, CardHeader } from '@mui/material';
import Image from 'next/image';
import rook from '$lib/image/rook.svg';

export default function EventCardBig() {
	return (
		<Card
			sx={{
				backgroundImage: 'url(/blekksprut2.jpg)',
				maxHeight: '267px',
				maxWidth: '306px',
				height: '100%',
				width: '100%',
				backgroundSize: 'cover',
				borderRadius: '1.75rem',
			}}
		>
			<CardHeader
				title="Tentakkel"
				titleTypographyProps={{ fontSize: '1.8rem' }}
				sx={{ height: '141px', alignItems: 'flex-end', padding: '1rem' }}
			/>
			<CardContent
				sx={{
					height: '126px',
					backgroundColor: 'rgba(0,0,0,0.5)',
					backdropFilter: 'blur(4px)',
					padding: '1rem',
				}}
			>
				<Typography sx={{ fontWeight: 'bold', fontSize: '1.1rem' }}> Gerhard Fajita </Typography>
				<Box
					sx={{
						display: 'flex',
						justifyContent: 'space-between',
						color: 'secondary.contrastText',
						paddingBottom: '0.5rem',
					}}
				>
					<Typography> Call of Cthulhu </Typography>
					<Box sx={{ display: 'flex', gap: '1rem' }}>
						<Box component={Image} priority src={rook} alt="rook icon" />
						<Box component={Image} priority src={rook} alt="rook icon" />
						<Box component={Image} priority src={rook} alt="rook icon" />
						<Box component={Image} priority src={rook} alt="rook icon" />
					</Box>
				</Box>
				<Typography sx={{ color: 'white' }}>
					Klarer du å overleve landets farligste fjell? Pass på at du ikke brenner deg!
				</Typography>
			</CardContent>
		</Card>
	);
}
