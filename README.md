# e-scooter

this is a demo project

e-scooter is a rent-a-ride model application that clients can intteract with to book a ride with mobile application or web
for simplicity we take an electric scooters

project has two components - web application and iot component, 'brains' of a scooter
there is thousands of this scooters around online (or not) that user can interract with

this scooter `fleet` is sending telemtry to application (kafka), like gps coordinates and other

entrypoint is a mobile app application (won't be implemented) 

the pipeline for this model is, in whole:

- authenticating in app
- make deposit\attach card
- finding a ride
- booking a ride
- riding
- releasing ride
- taking ride fare + returning deposit

flow diagrams https://github.com/integer00/diagrams
