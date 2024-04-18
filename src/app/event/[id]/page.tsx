'use client';
import { Box, Button, Chip, Icon, Paper, Slider, Typography } from '@mui/material';
import Image from 'next/image';
import NavigateBefore from '@mui/icons-material/NavigateBefore';
import IconButton from '@mui/material/IconButton';
import blekksprut2 from '$public/blekksprut2.jpg';
import HelpIcon from '@mui/icons-material/Help';
import { useState } from 'react';

const marks = [
	{
		value: 1,
		label: 'Ikke interessert',
	},
	{
		value: 2,
		label: 'Litt interessert',
	},
	{
		value: 3,
		label: 'Interessert',
	},
	{
		value: 4,
		label: 'Veldig interessert',
	},
];

const Event = () => {
	const arrayet = ['katt', 'hund', 'fugl', 'hatt', 'nisse'];
	const [interest, setInterest] = useState<Number>(0);
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
					placeholder="blur"
					loading="lazy"
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
				<Box
					sx={{
						backgroundColor: 'secondary.main',
						minHeight: '62px',
						textAlign: 'center',
						display: 'grid',
						placeContent: 'center',
						borderRadius: '0.2rem',
					}}
				>
					<Typography component="p">{marks[interest].label}</Typography>
				</Box>
				<Slider
					onChange={(e) => {
						const target = e.target as HTMLInputElement;
						setInterest(Number(target.value));
					}}
					aria-label="Temperature"
					defaultValue={0}
					valueLabelDisplay="off"
					shiftStep={1}
					step={1}
					min={0}
					max={3}
				/>
			</Box>
			<Box display="inline-flex" gap="0.4rem">
				<HelpIcon sx={{ scale: '1.5' }} />
				<Typography component="p">Forvirret? Les mer om påmeldingssystemet</Typography>
			</Box>

			<Typography component="p">
				Lorem ipsum dolor sit amet, consectetur adipisicing elit. Natus distinctio quia odit recusandae nobis
				autem, odio id pariatur magnam illo saepe laborum nesciunt quasi doloremque provident neque eligendi,
				quisquam quas?
			</Typography>
		</Box>
	);
};

export default Event;
