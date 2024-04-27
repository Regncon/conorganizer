'use client';
import { Box, Chip, Slider, Typography, chipClasses, sliderClasses, type SxProps, type Theme } from '@mui/material';
import Image from 'next/image';
import NavigateBefore from '@mui/icons-material/NavigateBefore';
import IconButton from '@mui/material/IconButton';
import blekksprut2 from '$public/blekksprut2.jpg';
import HelpIcon from '@mui/icons-material/Help';
import { useState } from 'react';
import Link from 'next/link';

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

	const paragraphStyle: SxProps<Theme> = { margin: '1rem 0' };
	const strongStyle: SxProps<Theme> = { fontWeight: 700 };
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

			<Box sx={{ display: 'flex', gap: '1rem', placeContent: 'center', marginBottom: 1 }}>
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
			<Box
				sx={{
					display: 'flex',
					gap: '.5em',
					overflowX: 'auto',
					paddingBottom: 1,
				}}
			>
				{arrayet.map((vesen) => (
					<Chip
						variant="outlined"
						label={vesen}
						key={vesen}
						sx={{
							[`&.${chipClasses.root} .${chipClasses.icon}, &`]: {
								color: 'secondary.contrastText',
								borderColor: 'secondary.contrastText',
							},
						}}
						icon={<NavigateBefore />}
					/>
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
				<Typography sx={paragraphStyle} component="p">
					{marks[interest].label}
				</Typography>
			</Box>
			<Box sx={{ padding: '0 0.35rem' }}>
				<Slider
					onChange={(e) => {
						const target = e.target as HTMLInputElement;
						setInterest(Number(target.value));
					}}
					sx={{
						color: 'secondary.contrastText',
						[`.${sliderClasses.rail}`]: {
							backgroundColor: '#3d3b3b',
							height: '1rem',
						},
						[`.${sliderClasses.track}`]: {
							height: '1rem',
						},
						[`.${sliderClasses.mark}`]: {
							height: '0.8rem',
							width: '0.8rem',
							borderRadius: '50%',
						},
						[`.${sliderClasses.markActive}`]: {
							backgroundColor: 'secondary.contrastText',
						},
					}}
					marks
					defaultValue={0}
					valueLabelDisplay="off"
					min={0}
					max={3}
				/>
			</Box>

			<Box
				component={Link}
				href="#"
				sx={{
					display: 'inline-flex',
					gap: '0.4rem',
					marginBottom: 2,
					paddingLeft: 1,
					color: 'secondary.contrastText',
				}}
			>
				<HelpIcon sx={{ scale: '1.5', placeSelf: 'center' }} />
				<Typography component="p">Forvirret? Les mer om påmeldingssystemet</Typography>
			</Box>
			<Typography sx={strongStyle} component="strong">
				Lorem ipsum dolor, sit amet consectetur
			</Typography>
			<Typography sx={strongStyle} component="strong">
				adipisicing elit. Nemo, illo quisquam. Quae odit impedit{' '}
			</Typography>
			<Typography sx={{ ...paragraphStyle, marginBottom: 0, paddingBlockEnd: '1rem' }} component="p">
				Lorem ipsum dolor sit amet, consectetur adipisicing elit. Natus distinctio quia odit recusandae nobis
				autem, odio id pariatur magnam illo saepe laborum nesciunt quasi doloremque provident neque eligendi,
				quisquam quas?
			</Typography>
		</Box>
	);
};

export default Event;
