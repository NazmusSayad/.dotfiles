#!/bin/bash

sudo spctl --master-disable

sudo mdutil -a -i off
sudo mdutil -a -E

sudo tmutil disablelocal
