# Changelog

[comment]: # (Changes since last release go here)

## 2.1.8 - Jun 27th 2023

* Added support for YAML serialization/deserialization

## 2.1.7 - Apr 13th 2023

* Renamed test_utils.go to utils_test.go

## 2.1.6 - Feb 15th 2023

* Added `GetAndMoveToBack()` and `GetAndMoveToFront()` methods

## 2.1.5 - Dec 13th 2022

* Added `Value()` method

## 2.1.4 - Dec 12th 2022

* Fixed a bug with UTF-8 special characters in JSON keys

## 2.1.3 - Dec 11th 2022

* Added support for JSON marshalling/unmarshalling of wrapper of primitive types

## 2.1.2 - Dec 10th 2022
* Allowing to pass options to `New`, to give a capacity hint, or initial data
* Allowing to deserialize nested ordered maps from JSON without having to explicitly instantiate them
* Added the `AddPairs` method

## 2.1.1 - Dec 9th 2022
* Fixing a bug with JSON marshalling

## 2.1.0 - Dec 7th 2022
* Added support for JSON serialization/deserialization
