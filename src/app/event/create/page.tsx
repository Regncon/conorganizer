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
} from '@mui/material';
import { Unstable_NumberInput as NumberInput } from '@mui/base/Unstable_NumberInput';
import CustomNumberInput from './CustomNumberInput';

const Create = () => {
	return (
		<form>
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
				<TextField name="system" label="Spillsystem" variant="outlined" />
			</Paper>
		</form>
	);
};
export default Create;
