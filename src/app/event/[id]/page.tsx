import { Box, Chip, Paper, Typography } from '@mui/material';
import Image from 'next/image';
import NavigateBefore from '@mui/icons-material/NavigateBefore';
import IconButton from '@mui/material/IconButton';
import blekksprut2 from '$public/blekksprut2.jpg';

const Event = () => {
	const arrayet = ['katt', 'hund', 'fugl', 'hatt', 'nisse'];
	return (
		<Box sx={{ display: 'grid', gridTemplateAreas: '"header""content"' }}>
			<Box
				sx={{
					gridArea: 'header',
					display: 'grid',
					'& > *': {
						gridColumn: ' 1 / 2',
						gridRow: ' 1 / 2',
					},
				}}
			>
				<Box
					component={Image}
					src={blekksprut2}
					alt="noe alt-tekst"
					sx={{
						width: '100%',
						maxWidth: '100%',
						aspectRatio: '3.3 / 2',
					}}
				/>
				<Box
					sx={{
						background: 'linear-gradient(0deg, black, transparent)',
					}}
				>
					<IconButton>
						<NavigateBefore />
					</IconButton>
					<Typography component="h1">Nei, du er en n00b!</Typography>
				</Box>
			</Box>
			<Box sx={{ gridArea: 'content' }}>
				<Box>
					<Typography component="h2">GM: Fransibald von Fokkoff</Typography>
					<Typography component="em" sx={{ color: 'orangered' }}>
						System: Mage - the ascension
					</Typography>
				</Box>
				<Box sx={{ display: 'flex', gap: '.5em' }}>
					{arrayet.map((vesen) => (
						<Chip variant="outlined" label={vesen} key={vesen} icon={<NavigateBefore />} />
					))}
				</Box>
				<Typography component="p">Lorem ipsum dolor, sit amet consectetur</Typography>
				<Typography component="p">adipisicing elit. Nemo, illo quisquam. Quae odit impedit </Typography>
				<Typography component="p">est, odio nisi doloremque ullam alias magnam aspernatur e</Typography>
				<Typography component="p">rror labore eligendi aliquid magni culpa neque ad.</Typography>
			</Box>
		</Box>
	);
};

export default Event;
