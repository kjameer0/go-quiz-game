# Describe the project

The basic idea with this project is that I am creating a timed quiz game based on the csv file provided by the exercise. The first level is creating the quiz without a timer, then creating it with the timer. I might be able to add more options in the long run as well.

## MVP

This app should allow a user to:

1. run the quiz
2. see a question
3. answer a question
4. see another question
5. answer questions until the list is exhausted
6. get a score

## What does the program need to be able to do?

1. read a csv file and grab questions and answers
2. count how many questions there are
3. write to stdout and take user input
4. check user input for correctness
5. mark a user as correct when they get an answer right

Relevant tools

1. flags package
2. csv package
3. os
4. goroutines
5. channels
6. timer package
