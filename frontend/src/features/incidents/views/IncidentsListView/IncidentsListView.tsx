import React from "react";
import { getIncidentsService } from "../../services/get-incidents.service";
import { IncidentsListClient } from "./IncidentsListClient";

export const IncidentsListView = async () => {
  let incidents;
  try {
    const res = await getIncidentsService();
    incidents = res.data;
  } catch (error) {
    throw error;
  }

  return (
    <IncidentsListClient initialIncidents={incidents} />
  );
};
