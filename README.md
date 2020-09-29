# SCOIR Technical Interview for Back-End Engineers
This repo contains an exercise intended for Back-End Engineers.

## Instructions
1. Fork this repo.
1. Using technology of your choice, complete [the assignment](./Assignment.md).
1. Update this README with
    * a `How-To` section containing any instructions needed to execute your program.
    * an `Assumptions` section containing documentation on any assumptions made while interpreting the requirements.
1. Before the deadline, submit a pull request with your solution.

## Expectations
1. Please take no more than 8 hours to work on this exercise. Complete as much as possible and then submit your solution.
1. This exercise is meant to showcase how you work. With consideration to the time limit, do your best to treat it like a production system.

## How To

### Requirements

- python 3.8+
- [pipenv](https://pypi.org/project/pipenv/)

### Dependencies

- [watchdog](https://pypi.org/project/watchdog/) for triggering actions when files are created

### Running

```
pipenv install; # INSTALL DEPENDENCIES
pipenv run python main.py; # EXECUTE PROGRAM WITH DEFAULT PARAMS
```

- [ ] TODO: create gif of execution
- [ ] TODO: document run command

## Assumptions

1. Data was generated via [mockaroo](https://www.mockaroo.com/a701ae50)
1. The program is intended to be used via Command Line Interface(CLI)
1. The program has been testing on MacOS and Linux(Ubuntu)
1. The program should run until the user exits the program
1. TODO: All the assumptions