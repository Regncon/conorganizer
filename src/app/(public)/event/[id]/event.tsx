'use client';
import { Box, Chip, Slider, Typography, sliderClasses, useTheme, type SxProps, type Theme } from '@mui/material';
import Image from 'next/image';
import NavigateBefore from '@mui/icons-material/NavigateBefore';
import IconButton from '@mui/material/IconButton';
import blekksprut2 from '$public/blekksprut2.jpg';
import HelpIcon from '@mui/icons-material/Help';
import { useState } from 'react';
import Link from 'next/link';
import { faChevronLeft, faUserSecret } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faScroll } from '@fortawesome/free-solid-svg-icons/faScroll';
import { Event } from '$lib/types';

export const dynamic = 'force-static';

const marks = [
	{
		value: 1,
		label: '🥱 Ikke interessert',
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

type props = {
	data: Event;
};
const MainEvent = ({ data }: props) => {
	const [interest, setInterest] = useState<number>(0);
	const {
		palette: {
			background: { default: themeBackgroundColor },
		},
	} = useTheme();

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
						background: `linear-gradient(0deg, ${themeBackgroundColor}, transparent)`,
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
							<FontAwesomeIcon icon={faChevronLeft} />
						</IconButton>
						<Typography
							variant="h1"
							align="center"
							sx={{ placeSelf: 'end center', paddingBottom: '2.5rem' }}
						>
							{data.title}
						</Typography>
					</Box>
				</Box>
			</Box>

			<Box sx={{ display: 'flex', gap: '1rem', marginBottom: '2rem' }}>
				<Box sx={{ display: 'flex', gap: '0.8rem', placeItems: 'center' }}>
					<Box component={FontAwesomeIcon} icon={faUserSecret} size="2x" sx={{ color: 'primary.main' }} />
					<Box>
						<Typography component="span" sx={{ color: 'primary.main' }}>
							{data.icons?.includes('rollespill') ? 'Gamemaster' : 'Arrangør'}
						</Typography>
						<Typography variant="h2">{data.gameMaster}</Typography>
					</Box>
				</Box>
				<Box sx={{ display: 'flex', gap: '0.8rem', placeItems: 'center' }}>
					<Box component={FontAwesomeIcon} icon={faScroll} size="2x" sx={{ color: 'primary.main' }} />
					<Box>
						<Typography component="span" sx={{ color: 'primary.main' }}>
							System
						</Typography>
						<Typography variant="h2">{data.system}</Typography>
					</Box>
				</Box>
			</Box>
			<Box
				sx={{
					display: 'flex',
					gap: '.5em',
					overflowX: 'auto',
					marginBottom: '4rem',
					paddingBottom: '0.35rem',
				}}
			>
				{data.icons?.map((vesen) => (
					<Chip label={vesen} key={vesen} color="primary" icon={<NavigateBefore />} />
				))}
			</Box>
			<Box
				sx={{
					backgroundColor: 'primary.main',
					color: 'primary.contrastText',
					minHeight: '62px',
					textAlign: 'center',
					display: 'grid',
					placeContent: 'center',
					borderRadius: '0.2rem',
					marginBottom: '1rem',
				}}
			>
				<Typography sx={paragraphStyle} component="p">
					{marks[interest].label}
				</Typography>
			</Box>
			<Box sx={{ padding: '0.35rem', marginBottom: '0.1rem' }}>
				<Slider
					onChange={(e) => {
						const target = e.target as HTMLInputElement;
						setInterest(Number(target.value));
					}}
					sx={{
						color: 'primary.main',
						[`.${sliderClasses.rail}`]: {
							backgroundColor: '#3d3b3b',
							height: '1rem',
						},
						[`.${sliderClasses.track}`]: {
							height: '1rem',
						},
						[`.${sliderClasses.mark}`]: {
							borderRadius: '50%',
							outlineColor: 'primary.main',
							outlineWidth: '0.8rem',
							outlineStyle: 'solid',
							outlineOffset: '-1px',
							opacity: '1',
						},
						[`.${sliderClasses.markActive}`]: {
							backgroundColor: 'primary.main',
						},
						[`.${sliderClasses.thumb}:before`]: {
							boxShadow: 'unset',
						},
						[`.${sliderClasses.valueLabelCircle}`]: {
							display: 'none',
						},
					}}
					marks
					defaultValue={0}
					min={0}
					max={3}
				/>
			</Box>

			<Typography sx={{ marginBottom: '0.8rem', textAlign: 'center' }}>
				Dra baren over for å melde din interesse!
			</Typography>

			<Box
				component={Link}
				href="#"
				sx={{
					display: 'inline-flex',
					gap: '0.4rem',
					marginBottom: '3rem',
					paddingLeft: '0.5rem',
					color: 'primary.main',
				}}
			>
				<HelpIcon sx={{ scale: '1.5', placeSelf: 'center' }} />
				<Typography component="p">Forvirret? Les mer om påmeldingssystemet</Typography>
			</Box>
			<Typography sx={strongStyle} component="strong">
				{data.shortDescription}
			</Typography>
			<Typography sx={{ ...paragraphStyle, marginBottom: 0, paddingBlockEnd: '1rem' }} component="p">
				Lorem ipsum dolor sit amet, consectetur adipisicing elit. Natus distinctio quia odit recusandae nobis
				autem, odio id pariatur magnam illo saepe laborum nesciunt quasi doloremque provident neque eligendi,
				quisquam quas?
			</Typography>
		</Box>
	);
};

export default MainEvent;
