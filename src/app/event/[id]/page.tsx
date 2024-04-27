'use client';
import { Box, Chip, Slider, Typography } from '@mui/material';
import Image from 'next/image';
import NavigateBefore from '@mui/icons-material/NavigateBefore';
import IconButton from '@mui/material/IconButton';
import blekksprut2 from '$public/blekksprut2.jpg';
import HelpIcon from '@mui/icons-material/Help';
import { useState } from 'react';

const marks = [
	{
		value: 1,
		label: '😡 Ikke interessert',
	},
	{
		value: 2,
		label: '😑 Litt interessert',
	},
	{
		value: 3,
		label: '😊 Interessert',
	},
	{
		value: 4,
		label: '🤩 Veldig interessert',
	},
];

const Event = () => {
	const arrayet = ['katt', 'hund', 'fugl', 'rollespill', 'nisse', 'nisse', 'nisse', 'nisse', 'nisse'];
	const [interest, setInterest] = useState<number>(0);
	return (
		<Box>
			<Box
				sx={{
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
					<Box
						sx={{
							display: 'grid',
							gridTemplateRows: '2rem 1fr',
							height: '100%',
							wordBreak: 'break-word',
						}}
					>
						<IconButton sx={{ placeSelf: 'start' }}>
							<NavigateBefore />
						</IconButton>
						<Typography
							variant="h1"
							align="center"
							sx={{ placeSelf: 'end center', paddingBottom: '2.5rem' }}
						>
							Nei, du er en n00b!
						</Typography>
					</Box>
				</Box>
			</Box>
			<Box>
				<Box sx={{ display: 'flex', gap: '1rem', placeContent: 'center' }}>
					<Box>
						<Typography component="span" sx={{ color: 'secondary.contrastText' }}>
							{arrayet.includes('rollespill') ? 'Gamemaster' : 'Arrangør'}
						</Typography>
						<Typography variant="h2">Fransibald von Fokkoff</Typography>
					</Box>
					<Box>
						<Typography component="span" sx={{ color: 'secondary.contrastText' }}>
							System
						</Typography>
						<Typography variant="h2">Mage - the ascension</Typography>
					</Box>
				</Box>
				<Box sx={{ display: 'flex', gap: '.5em', overflowX: 'auto' }}>
					{arrayet.map((vesen) => (
						<Chip variant="outlined" label={vesen} key={vesen} icon={<NavigateBefore />} />
					))}
				</Box>
				<Box
					sx={{
						backgroundColor: 'secondary.contrastText',
						minHeight: '62px',
						textAlign: 'center',
						display: 'grid',
						placeContent: 'center',
						borderRadius: '0.2rem',
					}}
				>
					<Typography component="p">{marks[interest].label}</Typography>
				</Box>
				<Box sx={{ padding: '0 0.35rem' }}>
					<Slider
						onChange={(e) => {
							const target = e.target as HTMLInputElement;
							setInterest(Number(target.value));
						}}
						aria-label="Temperature"
						marks
						defaultValue={0}
						valueLabelDisplay="off"
						shiftStep={1}
						step={1}
						min={0}
						max={3}
					/>
				</Box>
			</Box>
			<Box display="inline-flex" gap="0.4rem">
				<HelpIcon sx={{ scale: '1.5' }} />
				<Typography component="p">Forvirret? Les mer om påmeldingssystemet</Typography>
			</Box>
			<Typography component="strong">Lorem ipsum dolor, sit amet consectetur</Typography>
			<Typography component="strong">adipisicing elit. Nemo, illo quisquam. Quae odit impedit </Typography>
			<Typography component="p">
				Lorem ipsum dolor sit amet, consectetur adipisicing elit. Natus distinctio quia odit recusandae nobis
				autem, odio id pariatur magnam illo saepe laborum nesciunt quasi doloremque provident neque eligendi,
				quisquam quas?
			</Typography>
		</Box>
	);
};

export default Event;
