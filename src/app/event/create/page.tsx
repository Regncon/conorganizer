'use client';
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';
import TextField from '@mui/material/TextField';
import FormControl from '@mui/material/FormControl';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormLabel from '@mui/material/FormLabel';
import Radio from '@mui/material/Radio';
import RadioGroup from '@mui/material/RadioGroup';
import TextareaAutosize from '@mui/material/TextareaAutosize';
import FormGroup from '@mui/material/FormGroup';
import Checkbox from '@mui/material/Checkbox';
import Button from '@mui/material/Button';
import { Unstable_NumberInput as NumberInput } from '@mui/base/Unstable_NumberInput';
import Confetti from 'react-confetti';
import { useEffect, useState } from 'react';
import { doc, onSnapshot, setDoc } from 'firebase/firestore';
import { db } from '$lib/firebase/firebase';
import type { NewEvent } from '$app/types';
import Grid2 from '@mui/material/Unstable_Grid2/Grid2';

const Create = () => {
	const [isExploding, setIsExploding] = useState(false);
	const newEventDocRef = doc(db, 'users', 'HESPa0eFhPWWjGBGEXw6HjB1r8v1', 'myEvents', '72g29OPy3b9cCpmzg4R2');
	const [newEvent, setNewEvent] = useState<NewEvent>();
	console.log(newEvent);

	useEffect(() => {
		const unsubscribe = onSnapshot(newEventDocRef, (snapshot) => {
			setNewEvent(snapshot.data() as NewEvent);
			console.log(snapshot, 'snapshot');
		});

		return () => {
			unsubscribe();
		};
	}, []);

	return newEvent ?
			<Grid2
				container
				component="form"
                spacing='2rem'
				onChange={(e) => {
					const { value: inputValue, name: inputName, checked, type } = e.target as HTMLInputElement;

					let value: string | boolean = inputValue;
					let name: string = inputName;

					if (type === 'checkbox') {
						value = checked;
					}

					if (type === 'radio') {
						name = 'gameType';
						value = inputName;
					}

					setDoc(newEventDocRef, { ...newEvent, [name]: value });
				}}
			>
				{isExploding && (
					<Confetti
						onConfettiComplete={() => {
							setIsExploding(!isExploding);
						}}
						numberOfPieces={2000}
						recycle={false}
						height={document.body.scrollHeight}
						width={document.body.scrollWidth}
					/>
				)}
				<Paper>
					<Grid2 container gap="3rem">
						<Typography variant="h1">Meld på arrangement til Regncon XXXII 2024</Typography>
						<Typography>
							Takk for at du vil arrangere eit spel på Regncon, anten det er brettspel, kortspel,
							rollespel eller anna, så sett vi enormt pris på ditt bidrag. Fyll inn skjemaet så godt du
							kan, og ikkje vere redd for å ta kontakt med Regnconstyret på{' '}
							<a href="mailto:regncon@gmail.com">regncon@gmail.com</a>
							om du skulle lure på noko!
						</Typography>
					</Grid2>
				</Paper>

				<Grid2 xs={12}>
                    <Paper>
                        <TextField
                            name="title"
                            label="Tittel på spelmodul / arrangement"
                            value={newEvent.title}
                            variant="outlined"
                            required
                            fullWidth
                        />
                    </Paper>
				</Grid2>

				<Grid2 xs={12} md={6} lg={3}>
					<Paper>
						<TextField
							type="email"
							name="email"
							value={newEvent.email}
							label="E-postadresse"
							variant="outlined"
							required
							fullWidth
						/>
					</Paper>
				</Grid2>
				<Grid2 xs={12} md={6} lg={3}>
				<Paper>
					<TextField
						name="name"
						value={newEvent.name}
						label="Arrangørens namn (Ditt namn)"
						variant="outlined"
						required
						fullWidth
					/>
				</Paper>
				</Grid2>
				<Grid2 xs={12} md={6} lg={3}>
				<Paper>
					<TextField
						type="phone"
						name="phone"
						value={newEvent.phone}
						label="Kva telefonnummer kan vi nå deg på?"
						variant="outlined"
						required
						fullWidth
					/>
				</Paper>

				</Grid2>
				<Grid2 xs={12} md={6} lg={3}>
				<Paper>
					<TextField name="system" label="Spillsystem" value={newEvent.system} variant="outlined" fullWidth />
				</Paper>
				</Grid2>
								<Grid2 xs={12}>
				<Paper>
					<FormControl fullWidth>
						<FormLabel>Skildring av modulen (tekst til programmet):</FormLabel>
						<TextareaAutosize minRows={3} name="description" value={newEvent.description} fullWidth />
					</FormControl>
				</Paper>
				</Grid2>
				<Grid2 xs={12} md={4}>
				<Paper>
					<FormControl fullWidth>
						<FormLabel>Kva type spel er det?</FormLabel>
						<RadioGroup
							value={newEvent.gameType}
							aria-labelledby="demo-controlled-radio-buttons-group"
							name="controlled-radio-buttons-group"
						>
							<FormControlLabel
								value="rolePlaying"
								control={<Radio name="rolePlaying" />}
								label="rollespel"
							/>
							<FormControlLabel
								value="boardGame"
								control={<Radio name="boardGame" />}
								label="Brettspel"
							/>
							<FormControlLabel value="cardGame" control={<Radio name="cardGame" />} label="Kortspel" />
							<FormControlLabel value="other" control={<Radio name="other" />} label="Annet" />
						</RadioGroup>
					</FormControl>
				</Paper>
				</Grid2>
                <Grid2>
				<Paper>
				<Typography>Maks antall deltakere</Typography>
					<NumberInput
						value={newEvent.participants}
						slotProps={{
							input: {
								name: 'participants',
							},
						}}
					/>
				</Paper>
				</Grid2>

				<Paper>
					<FormGroup>
						<FormLabel>Kva for pulje kan du arrangere i?</FormLabel>
						<FormControlLabel
							control={<Checkbox checked={newEvent.fridayEvening} />}
							name="fridayEvening"
							label="Fredag Kveld"
						/>
						<FormControlLabel
							control={<Checkbox checked={newEvent.saturdayMorning} />}
							name="saturdayMorning"
							label="Lørdag Morgen"
						/>
						<FormControlLabel
							control={<Checkbox checked={newEvent.saturdayEvening} />}
							name="saturdayEvening"
							label="Lørdag Kveld"
						/>
						<FormControlLabel
							control={<Checkbox checked={newEvent.sundayMorning} />}
							name="sundayMorning"
							label="Søndag Morgen"
						/>
					</FormGroup>
				</Paper>
				<Paper>
					<FormGroup>
						<FormLabel>Kryss av for det som gjeld</FormLabel>
						<FormControlLabel
							control={<Checkbox name="moduleCompetition" checked={newEvent.moduleCompetition} />}
							label="Eg vil vere med på modulkonkurransen husk å sende modulen til moduler@regncon.no innen første september!)"
						/>
						<FormControlLabel
							control={<Checkbox name="childFriendly" checked={newEvent.childFriendly} />}
							label="Arrangementet passer for barn"
						/>
						<FormControlLabel
							control={<Checkbox name="adultsOnly" checked={newEvent.adultsOnly} />}
							label="Arrangementet passer berre for vaksne (18+)"
						/>
						<FormControlLabel
							control={<Checkbox name="beginnerFriendly" checked={newEvent.beginnerFriendly} />}
							label="Arrangementet er nybyrjarvenleg"
						/>
						<FormControlLabel
							control={<Checkbox name="possiblyEnglish" checked={newEvent.possiblyEnglish} />}
							label="Arrangementet kan haldast på engelsk"
						/>
						<FormControlLabel
							control={<Checkbox name="volunteersPossible" checked={newEvent.volunteersPossible} />}
							label="Andre kan halda arrangementet"
						/>
						<FormControlLabel
							control={<Checkbox name="lessThanThreeHours" checked={newEvent.lessThanThreeHours} />}
							label="Eg trur arrangementet vil vare kortare enn 3 timer"
						/>
						<FormControlLabel
							control={<Checkbox name="moreThanSixHours" checked={newEvent.moreThanSixHours} />}
							label="Eg trur arrangementet vil vare lenger enn 6 timer"
						/>
					</FormGroup>
				</Paper>
				<Paper>
					<FormControl>
						<FormLabel>Merknader: Er det noko anna du vil vi skal vite?</FormLabel>
						<TextareaAutosize minRows={3} name="additionalComments" value={newEvent.additionalComments} />
					</FormControl>
				</Paper>
				<Paper>
					<Typography>
						Skjemaet vert lagra automatisk, men om du likevel vil trykke på ein knapp, så er det ein her. :)
					</Typography>
					<Button onClick={() => setIsExploding(!isExploding)}>Send inn</Button>
				</Paper>
			</Grid2>
		:	null;
};
export default Create;

