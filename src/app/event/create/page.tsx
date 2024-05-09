'use client';
import {
	Box,
	Card,
	Typography,
	CardContent,
	Paper,
	Input,
	TextField,
	FormControl,
	FormControlLabel,
	FormLabel,
	Radio,
	RadioGroup,
	TextareaAutosize,
	FormGroup,
	Checkbox,
	Button,
} from '@mui/material';
import { Unstable_NumberInput as NumberInput } from '@mui/base/Unstable_NumberInput';
import CustomNumberInput from './CustomNumberInput';
import ConfettiExplosion from 'react-confetti-explosion';
import { useState } from 'react';

const Create = () => {
	const [isExploding, setIsExploding] = useState(false);
	return (
		<form>
			{isExploding && <ConfettiExplosion force={0.8} duration={3000} particleCount={250} width={1000} />}
			<Paper>
				<Typography variant="h1">Meld på arrangement til Regncon XXXII 2024</Typography>
				<Typography>
					Takk for at du vil arrangere eit spel på Regncon, anten det er brettspel, kortspel, rollespel eller
					anna, så sett vi enormt pris på ditt bidrag. Fyll inn skjemaet så godt du kan, og ikkje vere redd
					for å ta kontakt med Regnconstyret på <a href="mailto:regncon@gmail.com">regncon@gmail.com</a>
					om du skulle lure på noko!
				</Typography>
			</Paper>
			<Paper>
				<TextField type="email" name="email" label="E-postadresse" variant="outlined" required />
			</Paper>
			<Paper>
				<TextField name="name" label="Arrangørens namn (Ditt namn)" variant="outlined" required />
			</Paper>
			<Paper>
				<TextField
					type="phone"
					name="phone"
					label="Kva telefonnummer kan vi nå deg på?"
					variant="outlined"
					required
				/>
			</Paper>
			<Paper>
				<TextField
					name="title"
					label="Tittel på spelmodul / arrangement
"
					variant="outlined"
					required
				/>
			</Paper>
			<Paper>
				<TextField name="system" label="Spillsystem" variant="outlined" />
			</Paper>
			<Paper>
				<FormControl>
					<FormLabel>Kva type spel er det?</FormLabel>
					<RadioGroup defaultValue="rollespel" name="radio-buttons-group">
						<FormControlLabel name="rollespel" value="rollespel" control={<Radio />} label="Rollespel" />
						<FormControlLabel name="brettspel" value="brettspel" control={<Radio />} label="Brettspel" />
						<FormControlLabel name="kortspel" value="kortspel" control={<Radio />} label="Kortspel" />
						<FormControlLabel name="annet" value="annet" control={<Radio />} label="Annet" />
					</RadioGroup>
				</FormControl>
			</Paper>

			<Paper>
				<NumberInput />
			</Paper>
			<Paper>
				<FormControl>
					<FormLabel>Skildring av modulen (tekst til programmet):</FormLabel>
					<TextareaAutosize minRows={3} name="description" />
				</FormControl>
			</Paper>
			<Paper>
				<FormGroup>
					<FormLabel>Kva for pulje kan du arrangere i?</FormLabel>
					<FormControlLabel control={<Checkbox defaultChecked />} label="Fredag Kveld" />
					<FormControlLabel control={<Checkbox defaultChecked />} label="Lørdag Morgen" />
					<FormControlLabel control={<Checkbox defaultChecked />} label="Lørdag Kveld" />
					<FormControlLabel control={<Checkbox defaultChecked />} label="Søndag Morgen" />
				</FormGroup>
			</Paper>
			<Paper>
				<FormGroup>
					<FormLabel>Kryss av for det som gjeld</FormLabel>
					<FormControlLabel
						control={<Checkbox />}
						label="Eg vil vere med på modulkonkurransen husk å sende modulen til moduler@regncon.no innen første september!)"
					/>
					<FormControlLabel control={<Checkbox />} label="Arrangementet passer for barn" />
					<FormControlLabel control={<Checkbox />} label="Arrangementet passer berre for vaksne (18+)" />
					<FormControlLabel control={<Checkbox />} label="Arrangementet er nybyrjarvenleg" />
					<FormControlLabel control={<Checkbox />} label="Arrangementet kan haldast på engelsk" />
					<FormControlLabel control={<Checkbox />} label="Andre kan halda arrangementet" />
					<FormControlLabel
						control={<Checkbox />}
						label="Eg trur arrangementet vil vare kortare enn 3 timer"
					/>
					<FormControlLabel
						control={<Checkbox />}
						label="Eg trur arrangementet vil vare lenger enn 6 timer"
					/>
				</FormGroup>
			</Paper>
			<Paper>
				<FormControl>
					<FormLabel>Merknader: Er det noko anna du vil vi skal vite?</FormLabel>
					<TextareaAutosize minRows={3} name="description" />
				</FormControl>
			</Paper>
			<Paper>
				<Typography>
					Skjemaet vert lagra automatisk, men om du likevel vil trykke på ein knapp, så er det ein her. :)
				</Typography>
				<Button onClick={() => setIsExploding(!isExploding)}>Send inn</Button>
			</Paper>
		</form>
	);
};
export default Create;
