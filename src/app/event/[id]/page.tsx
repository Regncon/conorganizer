import { Box, Chip, Typography } from '@mui/material';
import Image from 'next/image';
import NavigateBefore from '@mui/icons-material/NavigateBefore';
import IconButton from '@mui/material/IconButton';

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
				<Image src="/blekksprut2.jpg" alt="noe alt-tekst" width={375} height={230} />
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
				<Typography component="p">
					Lorem ipsum dolor, sit amet consectetur adipisicing elit. Nemo, illo quisquam. Quae odit impedit
					est, odio nisi doloremque ullam alias magnam aspernatur error labore eligendi aliquid magni culpa
					neque ad.
				</Typography>
			</Box>
		</Box>
	);
};

export default Event;
