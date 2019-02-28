# 2SMS
SCION Monitoring System

## Project Structure
- common:       contains modules used in more than one component
- endpoint:     contains all files related to the Endpoint component
- manager:      contains all files related to the Manager component
- scraper:      contains all files related to the Scraper component
- storage:      contains all files related to the Storage component
- services:     
- alerting:     

## Building from source
Make sure that you have ```go``` (at least version 1.9) and `SCION` () installed. Then all you need to is is running ```go build <component>.go``` in the corresponding ```<component>``` folder. E.g. to build the manager component we have to run ```go build manager.go``` in ```/manager```.
