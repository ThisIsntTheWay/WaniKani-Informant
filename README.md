# WaniKani informant
This little tool informs you if graduating reviews are coming up using discord webhooks.
It is built with Go 1.18+, but might be backwards-compatible with prior versions.

A "graduating" review is any review that contains items currently sitting at SRS stage 4.  
If these items were reviewed successfully, they'd advance to stage "Guru", thus meaning you've 'passed' them.

## External modules
The following external modules are used:
 - github.com/TwiN/go-color
   - Colored console output
 - github.com/go-resty/resty/v2
   - REST client

### Preview
The message will show color-coded information for each item category with a graduation chance:
![image](https://user-images.githubusercontent.com/13659371/172054625-85259e01-3de2-4c2b-af6e-78d0b8d1c1ed.png)

The webhooks message can be customized in `json/msgGraduationTemplate.json`.