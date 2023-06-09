#!/bin/bash
psql --dbname "$POSTGRES_DB" --username "$POSTGRES_USER" -c 'create database class_booking_qa with owner class_booking;';
