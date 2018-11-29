FROM jupyter/minimal-notebook

ENV MCSTARTUP $HOME/.ipython/profile_default/startup
RUN mkdir -p $MCSTARTUP
COPY mysql_magic.py $MCSTARTUP

